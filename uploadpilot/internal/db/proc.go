package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProcessorRepo struct {
	collectionName string
}

func NewProcessorRepo() *ProcessorRepo {
	return &ProcessorRepo{
		collectionName: "processors",
	}
}

func (i *ProcessorRepo) GetAll(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	collection := db.Collection(i.collectionName)
	var processors []models.Processor
	opts := options.Find().SetProjection(bson.M{"tasks": 0}).SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"workspaceId": wsID}, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &processors); err != nil {
		return nil, err
	}

	return processors, nil
}

func (i *ProcessorRepo) Get(ctx context.Context, workspaceID string, processorID string) (*models.Processor, error) {
	id, err := primitive.ObjectIDFromHex(processorID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, processorID)
	}

	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	collection := db.Collection(i.collectionName)
	var cb models.Processor
	err = collection.FindOne(ctx, bson.M{"workspaceId": wsID, "_id": id}).Decode(&cb)
	if err != nil {
		return nil, err
	}

	return &cb, nil
}

func (i *ProcessorRepo) Create(ctx context.Context, workspaceID string, processor *models.Processor) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}
	now := primitive.NewDateTimeFromTime(time.Now())
	processor.WorkspaceID = wsID
	processor.ID = primitive.NewObjectID()
	processor.UpdatedAt = now
	processor.UpdatedBy = user.Email
	processor.CreatedAt = now
	processor.CreatedBy = user.Email
	processor.Enabled = true

	// validate before insert
	if err := infra.Validator.ValidateBody(processor); err != nil {
		return err
	}

	collection := db.Collection(i.collectionName)
	_, err = collection.InsertOne(ctx, processor)
	if err != nil {
		infra.Log.Errorf("failed to create processor: %s", err.Error())
		return err
	}

	return nil
}

func (wr *ProcessorRepo) Patch(ctx context.Context, workspaceID, processorID string, patch bson.M) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	id, err := primitive.ObjectIDFromHex(processorID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, processorID)
	}

	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}

	patch["updatedBy"] = user.Email
	patch["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	collection := db.Collection(wr.collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"workspaceId": wsID, "_id": id}, bson.M{"$set": patch})
	return err
}

func (wr *ProcessorRepo) Delete(ctx context.Context, workspaceID string, processorID string) error {
	wsID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, workspaceID)
	}

	id, err := primitive.ObjectIDFromHex(processorID)
	if err != nil {
		return fmt.Errorf(msg.InvalidObjectID, processorID)
	}

	collection := db.Collection(wr.collectionName)
	if _, err := collection.DeleteOne(ctx, bson.M{"workspaceId": wsID, "_id": id}); err != nil {
		return err
	}

	return nil
}
