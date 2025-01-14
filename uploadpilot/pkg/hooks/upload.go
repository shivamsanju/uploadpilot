package hooks

import (
	"fmt"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/shivamsanju/uploader/internal/db/models"
	"github.com/shivamsanju/uploader/internal/db/repo"
	"github.com/shivamsanju/uploader/pkg/cloudstorage"
	"github.com/shivamsanju/uploader/pkg/globals"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadToDatastoreHook(hook tusd.HookEvent, imp *models.Import) error {
	headers := hook.HTTPRequest.Header
	uploaderId := headers.Get("uploaderId")
	if len(uploaderId) == 0 {
		return fmt.Errorf("missing uploaderId in header")
	}
	uploaderIDObjID, err := primitive.ObjectIDFromHex(uploaderId)
	if err != nil {
		return fmt.Errorf("invalid uploaderId: %w", err)
	}
	imp.UploaderID = uploaderIDObjID
	uploaderRepo := repo.NewUploaderRepo()
	uploader, err := uploaderRepo.GetDataStoreCreds(hook.Context, uploaderId)
	if err != nil {
		return err
	}
	if uploader == nil {
		return fmt.Errorf("uploader not found")
	}
	bucket, ok := uploader["bucket"].(string)
	if !ok {
		return fmt.Errorf("bucket not found")
	}
	connectorDetails, ok := uploader["connectorDetails"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("connectorDetails not found")
	}

	var storageConnectorDetails models.StorageConnector
	if err := mapstructure.Decode(connectorDetails, &storageConnectorDetails); err != nil {
		return fmt.Errorf("failed to map connectorDetails to StorageConnector: %w", err)
	}

	cloudUploader, err := cloudstorage.NewUploader(&storageConnectorDetails, bucket)
	if err != nil {
		return err
	}
	uploadPath := path.Join(globals.TusUploadDir, hook.Upload.ID)
	objectFileName := hook.Upload.ID + "_" + hook.Upload.MetaData["filename"]
	imp.StoredFileName = objectFileName
	err = cloudUploader.Upload(uploadPath, objectFileName)
	if err != nil {
		return err
	}
	return nil
}
