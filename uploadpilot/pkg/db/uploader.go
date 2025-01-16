package db

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UploaderRepo interface {
	GetAll(ctx context.Context) ([]models.Uploader, error)
	Get(ctx context.Context, id string) (*bson.M, error)
	Create(ctx context.Context, cb *models.Uploader) (primitive.ObjectID, error)
	Delete(ctx context.Context, id string) error
	GetDataStoreCreds(ctx context.Context, id string) (map[string]interface{}, error)
	UpdateConfig(ctx context.Context, id string, cb *models.UploaderConfig, updatedBy string) error
}

type uploaderRepo struct {
	collectionName string
}

func NewUploaderRepo() UploaderRepo {
	return &uploaderRepo{
		collectionName: "uploaders",
	}
}

func (ur *uploaderRepo) GetAll(ctx context.Context) ([]models.Uploader, error) {
	collection := g.Db.Database(g.DbName).Collection(ur.collectionName)
	var cb []models.Uploader
	opts := options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		g.Log.Errorf("no Uploaders found: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &cb)
	g.Log.Infof("found %d Uploaders", len(cb))
	return cb, nil
}

func (ur *uploaderRepo) Get(ctx context.Context, id string) (*bson.M, error) {
	collection := g.Db.Database(g.DbName).Collection(ur.collectionName)
	UploaderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id": UploaderID,
		}},
		{"$lookup": bson.M{
			"from":         "datastores",
			"localField":   "dataStoreId",
			"foreignField": "_id",
			"as":           "dataStoreDetails",
		}},
		{"$unwind": bson.M{"path": "$dataStoreDetails"}},
		{"$lookup": bson.M{
			"from":         "storageconnectors",
			"localField":   "dataStoreDetails.connectorId",
			"foreignField": "_id",
			"as":           "dataStoreDetails.connectorDetails",
		}},
		{"$unwind": bson.M{"path": "$dataStoreDetails.connectorDetails"}},
		{"$addFields": bson.M{
			"dataStoreDetails.connectorName": "$dataStoreDetails.connectorDetails.name",
			"dataStoreDetails.connectorType": "$dataStoreDetails.connectorDetails.type",
			"dataStoreDetails.connectorId":   "$dataStoreDetails.connectorDetails._id",
		}},
		{"$unset": "dataStoreDetails.connectorDetails"},
	}
	cb := []bson.M{}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		g.Log.Errorf("failed to find Uploader: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to find Uploader: %s", err.Error())
		return nil, err
	}
	g.Log.Infof("found Uploader: %+v", cb)
	if len(cb) == 0 {
		return nil, nil
	}
	return &cb[0], nil
}

func (ur *uploaderRepo) Create(ctx context.Context, cb *models.Uploader) (primitive.ObjectID, error) {
	cb.ID = primitive.NewObjectID()
	cb.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	cb.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := g.Db.Database(g.DbName).Collection(ur.collectionName)
	r, err := collection.InsertOne(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to add Uploader: %v", err.Error())
		return primitive.ObjectID{}, err
	}
	return (r.InsertedID).(primitive.ObjectID), nil
}

func (ur *uploaderRepo) Delete(ctx context.Context, id string) error {
	collection := g.Db.Database(g.DbName).Collection(ur.collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (ur *uploaderRepo) UpdateConfig(ctx context.Context, id string, updatedData *models.UploaderConfig, updatedBy string) error {
	updatedAt := primitive.NewDateTimeFromTime(time.Now())
	uploaderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	collection := g.Db.Database(g.DbName).Collection(ur.collectionName)
	g.Log.Infof("\n\n\ns -----------> %+v\n\n\n", updatedData)
	update := bson.M{
		"$set": bson.M{
			"config.maxFileSize":            updatedData.MaxFileSize,
			"config.minFileSize":            updatedData.MinFileSize,
			"config.maxNumberOfFiles":       updatedData.MaxNumberOfFiles,
			"config.minNumberOfFiles":       updatedData.MinNumberOfFiles,
			"config.maxTotalFileSize":       updatedData.MaxTotalFileSize,
			"config.allowedFileTypes":       updatedData.AllowedFileTypes,
			"config.allowedSources":         updatedData.AllowedSources,
			"config.requiredMetadataFields": updatedData.RequiredMetadataFields,
			"config.theme":                  updatedData.Theme,
			"config.showStatusBar":          updatedData.ShowStatusBar,
			"config.showProgress":           updatedData.ShowProgress,
			"config.allowPauseAndResume":    updatedData.AllowPauseAndResume,
			"config.enableImageEditing":     updatedData.EnableImageEditing,
			"config.useCompression":         updatedData.UseCompression,
			"config.useFaultTolerantMode":   updatedData.UseFaultTolerantMode,
			"updatedBy":                     updatedBy,
			"updatedAt":                     updatedAt,
		},
	}
	g.Log.Infof("ddd %+v", update)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": uploaderID}, update)
	return err
}

func (ur *uploaderRepo) GetDataStoreCreds(ctx context.Context, id string) (map[string]interface{}, error) {
	collection := g.Db.Database(g.DbName).Collection(ur.collectionName)
	uploaderID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id": uploaderID,
		}},
		{"$lookup": bson.M{
			"from":         "datastores",
			"localField":   "dataStoreId",
			"foreignField": "_id",
			"as":           "dataStoreDetails",
		}},
		{"$unwind": bson.M{"path": "$dataStoreDetails"}},
		{"$lookup": bson.M{
			"from":         "storageconnectors",
			"localField":   "dataStoreDetails.connectorId",
			"foreignField": "_id",
			"as":           "connectorDetails",
		}},
		{"$unwind": bson.M{"path": "$connectorDetails"}},
		{"$addFields": bson.M{
			"bucket": "$dataStoreDetails.bucket",
		}},
		{"$unset": "dataStoreDetails"},
	}
	cb := []map[string]interface{}{}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		g.Log.Errorf("failed to find Uploader: %s", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &cb)
	if err != nil {
		g.Log.Errorf("failed to find Uploader: %s", err.Error())
		return nil, err
	}
	g.Log.Infof("found Uploader: %+v", cb)
	if len(cb) == 0 {
		return nil, nil
	}
	return cb[0], nil
}
