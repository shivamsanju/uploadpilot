package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

func DownloadFileFromS3(ctx context.Context, objectKey, localFilePath string) error {
	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer file.Close()

	obj, err := infra.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &config.S3BucketName,
		Key:    &objectKey,
	})
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer obj.Body.Close()

	if _, err = io.Copy(file, obj.Body); err != nil {
		return fmt.Errorf("failed to copy S3 object to temporary file: %w", err)
	}

	return nil
}

func DownloadAndExtractFileFromS3(ctx context.Context, objectKey, localDir string) error {
	tmpFilePath := filepath.Join(os.TempDir(), objectKey+".zip")
	defer os.RemoveAll(tmpFilePath)
	if err := DownloadFileFromS3(ctx, objectKey, tmpFilePath); err != nil {
		return fmt.Errorf("failed to download file from S3: %w", err)
	}

	if err := UnzipFile(ctx, tmpFilePath, localDir); err != nil {
		return fmt.Errorf("failed to unzip file: %w", err)
	}

	return nil
}

// ZipAndUpload cleaning up dirToZip is users task
func ZipAndUploadToS3(ctx context.Context, dirToZip, zipFileName string) (string, error) {
	objectName := zipFileName + ".zip"
	tmpDir := os.TempDir()
	zipFilePath := filepath.Join(tmpDir, objectName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}

	defer func() {
		zipFile.Close()
		os.RemoveAll(tmpDir)
	}()

	if err := ZipDirectory(ctx, dirToZip, zipFile); err != nil {
		return "", fmt.Errorf("failed to zip extracted content: %w", err)
	}

	zipFile.Seek(0, io.SeekStart)

	_, err = infra.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &config.S3BucketName,
		Key:    &objectName,
		Body:   zipFile,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload zip file to S3: %w", err)
	}
	return objectName, nil
}
