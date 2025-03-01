package data

import (
	"context"
)

type ActivityDataContainer interface {
	GetContainerID() string
	DownloadFile(ctx context.Context, downloadPath string) error
	DownloadAndUnzip(ctx context.Context, downloadDir string) error
	UploadFile(ctx context.Context, filepath string) error
	ZipAndUploadDirectory(ctx context.Context, sourceDir string) error
}
