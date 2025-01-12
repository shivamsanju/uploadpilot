package cloudstorage

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/chartmuseum/storage"
	"github.com/shivamsanju/uploader/internal/db/models"
)

type Uploader struct {
	Backend storage.Backend
}

func NewUploader(connector *models.StorageConnector, bucket string) (*Uploader, error) {
	if connector.Type == models.StorageTypeS3 {
		backend := storage.NewAmazonS3BackendWithCredentials(
			bucket,
			"",
			connector.S3Config.Region,
			"",
			"",
			credentials.NewStaticCredentials(connector.S3Config.AccessKey, connector.S3Config.SecretKey, ""),
		)
		uploader := Uploader{Backend: backend}
		return &uploader, nil
	}
	return nil, fmt.Errorf("unsupported storage type: %s", connector.Type)
}

func (uploader *Uploader) Upload(filepath, filename string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = uploader.Backend.PutObject(filename, content)
	if err != nil {
		return err
	}
	fmt.Printf("%s successfully uploaded\n", filename)
	return nil
}
