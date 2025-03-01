package data

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/agent/internal/zip"
)

type s3DataContainer struct {
	client *s3.Client
	bucket string
	key    string
}

func NewS3DataContainer(client *s3.Client, bucket, key string) ActivityDataContainer {
	return &s3DataContainer{
		client: client,
		bucket: bucket,
		key:    key,
	}
}

func (s *s3DataContainer) GetContainerID() string {
	return s.key
}

// DownloadFile downloads a file from S3 and saves it locally.
func (s *s3DataContainer) DownloadFile(ctx context.Context, downloadPath string) error {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	})
	if err != nil {
		return fmt.Errorf("failed to download file from S3: %w", err)
	}
	defer output.Body.Close()

	file, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, output.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *s3DataContainer) DownloadAndUnzip(ctx context.Context, downloadDir string) error {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	})
	if err != nil {
		return fmt.Errorf("failed to download zip file: %w", err)
	}
	defer output.Body.Close()

	tmpfile, err := os.CreateTemp("", "momentum-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, output.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return zip.UnzipFile(ctx, tmpfile.Name(), downloadDir)
}

func (s *s3DataContainer) UploadFile(ctx context.Context, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *s3DataContainer) ZipAndUploadDirectory(ctx context.Context, sourceDir string) error {
	var buf bytes.Buffer
	if err := zip.ZipDirectory(ctx, sourceDir, &buf); err != nil {
		return fmt.Errorf("failed to zip directory: %w", err)
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(s.key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("application/zip"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload zip file: %w", err)
	}

	return nil
}
