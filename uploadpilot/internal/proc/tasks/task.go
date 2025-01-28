package tasks

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskKey string

type Task interface {
	Do(ctx context.Context) error
	AddInput(data *TaskData)
	GetOutput() *TaskData
}

type TaskData struct {
	RunID           string
	ProcessorID     string
	WorkspaceID     string
	UploadID        string
	TmpDir          string
	TmpDirLock      *sync.Mutex
	ContinueOnError bool
}

func NewTaskData(workspaceID, processorID, uploadID, tmpDir string, continueOnError bool) *TaskData {
	return &TaskData{
		RunID:           primitive.NewObjectID().Hex(),
		ProcessorID:     processorID,
		WorkspaceID:     workspaceID,
		UploadID:        uploadID,
		TmpDir:          tmpDir,
		ContinueOnError: continueOnError,
		TmpDirLock:      &sync.Mutex{},
	}
}

type BaseTask struct {
	Data *TaskData
}

func NewBaseTask() *BaseTask {
	return &BaseTask{
		Data: &TaskData{},
	}
}

func (t *BaseTask) AddInput(data *TaskData) {
	t.Data = data
}

func (t *BaseTask) GetOutput() *TaskData {
	return t.Data
}
