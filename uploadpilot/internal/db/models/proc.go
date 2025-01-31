package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskKey string

type ProcTask struct {
	ID              string  `bson:"id" json:"id" validate:"required"`
	Key             TaskKey `bson:"key" json:"key" validate:"required"`
	Type            string  `bson:"type" json:"type" validate:"required"`
	Retry           uint64  `bson:"retry" json:"retry"`
	ContinueOnError bool    `bson:"continueOnError" json:"continueOnError"`
	TimeoutMilSec   uint64  `bson:"timeoutMilSec" json:"timeoutMilSec"`
	Data            JSON    `bson:"data" json:"data" validate:"required"`
	Position        JSON    `bson:"position" json:"position"`
	Measured        JSON    `bson:"measured" json:"measured"`
	Deletable       bool    `bson:"deletable" json:"deletable"`
}

type ProcTaskEdge struct {
	ID        string `bson:"id" json:"id" validate:"required"`
	Source    string `bson:"source" json:"source" validate:"required"`
	Target    string `bson:"target" json:"target" validate:"required"`
	Deletable bool   `bson:"deletable" json:"deletable"`
}

type ProcTaskCanvas struct {
	Nodes []ProcTask     `bson:"nodes" json:"nodes" validate:"required"`
	Edges []ProcTaskEdge `bson:"edges" json:"edges" validate:"required"`
}

type Processor struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	WorkspaceID primitive.ObjectID `bson:"workspaceId" json:"workspaceId" validate:"required"`
	Triggers    []string           `bson:"triggers" json:"triggers" validate:"required"`
	Tasks       *ProcTaskCanvas    `bson:"tasks" json:"tasks" validate:"required"`
	Enabled     bool               `bson:"enabled" json:"enabled"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt   primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}
