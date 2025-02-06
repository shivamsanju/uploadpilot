package tasks

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

func (t *BaseTask) GetTaskOutDir() string {
	return t.TmpDir + "/" + t.TaskID + "/out"
}

func (t *BaseTask) GetTaskInputDir() string {
	return t.TmpDir + "/" + t.TaskID + "/raw"
}

func (t *BaseTask) GetFileTypeFromS3(ctx context.Context, objectKey *string) (*string, error) {
	infra.Log.Infof("objectKey: %s checking type from s3", *objectKey)
	file, err := infra.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &config.S3BucketName,
		Key:    objectKey,
	})
	if err != nil {
		infra.Log.Infof("objectKey: %s could not check type from s3. error: %s", *objectKey, err.Error())

		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	return file.ContentType, nil
}

func (t *BaseTask) SaveInputFile(ctx context.Context) error {
	infra.Log.Infof("inputObjId: %+v, taskID: %s", t.Input, t.TaskID)
	inputObjId, ok := t.Input["inputObjId"].(string)
	if !ok {
		return fmt.Errorf("input inputObjId is not a string")
	}

	taskInputDir := t.GetTaskInputDir()
	os.MkdirAll(taskInputDir, os.ModePerm)

	ct, err := t.GetFileTypeFromS3(ctx, &inputObjId)
	if err != nil {
		return err
	}
	infra.Log.Infof("file type: %s", *ct)

	if *ct == "application/zip" {
		return utils.DownloadAndExtractFileFromS3(ctx, inputObjId, taskInputDir)

	}
	return utils.DownloadFileFromS3(ctx, inputObjId, taskInputDir+"/"+inputObjId)
}

func (t *BaseTask) SaveOutputFile(ctx context.Context) (string, error) {
	objectName, err := utils.ZipAndUploadToS3(ctx, t.GetTaskOutDir(), t.TaskID)
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func (t *BaseTask) Setup() error {
	inpDir := t.GetTaskInputDir()
	outDir := t.GetTaskOutDir()

	if err := os.MkdirAll(inpDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (t *BaseTask) Cleanup() error {
	return os.RemoveAll(t.TmpDir + "/" + t.TaskID)
}
