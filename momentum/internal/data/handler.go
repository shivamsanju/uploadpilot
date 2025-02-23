package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/phuslu/log"
	"github.com/uploadpilot/go-core/dsl"
	"github.com/uploadpilot/momentum/internal/config"
	"github.com/uploadpilot/momentum/internal/msg"
)

type ActivityDataHandler struct {
	wfMetaStr        string
	WorkflowMeta     *dsl.WorkflowMeta
	activityKey      string
	inputActivityKey string
	inputDC          ActivityDataContainer
	selfDC           ActivityDataContainer
	InputDir         string
	OutputDir        string
	RunInfo          RunInfo
	s3Client         *s3.Client
	s3Bucket         string
}

func NewActivityDataHandler(ctx context.Context, wfMeta string, activityKey, inputActivityKey string, s3Client *s3.Client) *ActivityDataHandler {
	return &ActivityDataHandler{
		wfMetaStr:        wfMeta,
		inputActivityKey: inputActivityKey,
		activityKey:      activityKey,
		s3Client:         s3Client,
		s3Bucket:         config.GetAppConfig().S3BucketName,
	}
}

func (a *ActivityDataHandler) PrepareDataLayer(ctx context.Context) error {
	if err := a.LoadRunInfo(ctx); err != nil {
		log.Error().Err(err).Msg("failed to load run info")
		return errors.New(msg.ErrRunInfoFileNotFound)
	}

	selfDataContainerID := uuid.New().String()
	a.selfDC = NewS3DataContainer(a.s3Client, a.s3Bucket, selfDataContainerID)

	inputDataContainerID := a.WorkflowMeta.UploadID
	if a.inputActivityKey != "" {
		inputActivityInfo, ok := a.RunInfo.ActivityInfoMap[a.inputActivityKey]
		if !ok {
			log.Error().Msg("failed to find input activity info")
			return errors.New(msg.ErrTaskInfoNotFound)
		}
		inputDataContainerID = inputActivityInfo.DataContainerID
	}

	a.inputDC = NewS3DataContainer(a.s3Client, a.s3Bucket, inputDataContainerID)

	a.InputDir = os.TempDir() + "/" + uuid.NewString()
	if err := os.Mkdir(a.InputDir, 0755); err != nil {
		log.Error().Err(err).Msg("failed to create input temp dir")
		return errors.New(msg.ErrTmpDirCreationFailed)
	}

	a.OutputDir = os.TempDir() + "/" + selfDataContainerID
	if err := os.Mkdir(a.OutputDir, 0755); err != nil {
		log.Error().Err(err).Msg("failed to create output temp dir")
		return errors.New(msg.ErrTmpDirCreationFailed)
	}

	if inputDataContainerID == a.WorkflowMeta.UploadID {
		dp := filepath.Join(a.InputDir, a.WorkflowMeta.UploadFileName)
		if err := a.inputDC.DownloadFile(ctx, dp); err != nil {
			log.Error().Err(err).Msg("failed to download input file")
			return errors.New(msg.ErrInputDownloadFailed)
		}
	} else {
		if err := a.inputDC.DownloadAndUnzip(ctx, a.InputDir); err != nil {
			log.Error().Err(err).Msg("failed to download and unzip input file")
			return errors.New(msg.ErrInputDownloadFailed)
		}
	}

	return nil
}

func (a *ActivityDataHandler) SaveOutput(ctx context.Context) (string, error) {
	numfiles := 0
	numBytes := int64(0)
	files := []FileInfo{}

	err := filepath.Walk(a.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		numfiles++
		numBytes += info.Size()

		files = append(files, FileInfo{
			Name: info.Name(),
			Size: info.Size(),
			Ext:  filepath.Ext(info.Name()),
		})
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to walk output dir")
		return "", errors.New(msg.ErrOutputSaveFailed)
	}

	a.RunInfo.ActivityInfoMap[a.activityKey] = ActivityInfo{
		ActivityKey:     a.activityKey,
		DataContainerID: a.selfDC.GetContainerID(),
		NumFiles:        numfiles,
		NumBytes:        numBytes,
		Files:           files,
	}

	if err := a.SaveRunInfo(ctx); err != nil {
		log.Error().Err(err).Msg("failed to save run info")
		return "", errors.New(msg.ErrRunInfoSaveFailed)
	}

	if err := a.selfDC.ZipAndUploadDirectory(ctx, a.OutputDir); err != nil {
		log.Error().Err(err).Msg("failed to zip and save output dir")
		return "", errors.New(msg.ErrOutputSaveFailed)
	}

	return a.selfDC.GetContainerID(), nil
}

func (a *ActivityDataHandler) Cleanup(ctx context.Context) {
	err := os.RemoveAll(a.InputDir)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove input dir")
	}

	err = os.RemoveAll(a.OutputDir)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove output dir")
	}
}

func (a *ActivityDataHandler) LoadRunInfo(ctx context.Context) error {
	var wfmeta dsl.WorkflowMeta
	if err := json.Unmarshal([]byte(a.wfMetaStr), &wfmeta); err != nil {
		return errors.New(msg.ErrInvalidWorkflowMetadata)
	}

	log.Info().Msgf("workflow metadata: %+v", wfmeta)
	a.WorkflowMeta = &wfmeta
	infoFile, err := a.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(config.GetAppConfig().S3BucketName),
		Key:    aws.String(a.WorkflowMeta.RunID + ".info"),
	})

	if err != nil {
		var apiErr *types.NoSuchKey
		if errors.As(err, &apiErr) {
			a.RunInfo = RunInfo{
				ActivityInfoMap: map[string]ActivityInfo{},
				RunID:           a.WorkflowMeta.RunID,
				UploadID:        a.WorkflowMeta.UploadID,
				WorkflowID:      a.WorkflowMeta.WorkflowID,
				ProcessorID:     a.WorkflowMeta.ProcessorID,
				WorkspaceID:     a.WorkflowMeta.WorkspaceID,
			}
			return nil
		}
		log.Error().Err(err).Msg("failed to get run info file")
		return errors.New(msg.ErrRunInfoFileNotFound)
	}

	a.RunInfo = RunInfo{}
	defer infoFile.Body.Close()
	infoBytes, err := io.ReadAll(infoFile.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read run info file")
		return errors.New(msg.ErrRunInfoFileReadFailed)
	}

	if err := json.Unmarshal(infoBytes, &a.RunInfo); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal run info file")
		return errors.New(msg.ErrRunInfoFileUnmarshalFailed)
	}

	return nil
}

func (a *ActivityDataHandler) SaveRunInfo(ctx context.Context) error {
	infoBytes, err := json.Marshal(a.RunInfo)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal run info file")
		return errors.New(msg.ErrRunInfoFileMarshalFailed)
	}

	_, err = a.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.GetAppConfig().S3BucketName),
		Key:    aws.String(a.WorkflowMeta.RunID + ".info"),
		Body:   bytes.NewReader(infoBytes),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to save run info file")
		return errors.New(msg.ErrRunInfoFileSaveFailed)
	}

	return nil
}

func (a *ActivityDataHandler) HandleActivity(ctx context.Context) error {
	return nil
}
