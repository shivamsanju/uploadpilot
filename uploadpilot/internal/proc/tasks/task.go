package tasks

import (
	"context"
	"sync"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskKey string

type Task interface {
	Do(ctx context.Context) error
	GetTask() *BaseTask
	MakeTask(data *BaseTask)
}

type BaseTask struct {
	RunID       string
	TaskID      string
	ProcessorID string
	WorkspaceID string
	UploadID    string
	TmpDir      string
	TmpDirLock  *sync.Mutex
	TaskParams  models.JSON
	Input       map[string]interface{}
	Output      map[string]interface{}
}

func NewBaseTask(workspaceID, processorID, uploadID, tmpDir, inputObjId string) *BaseTask {
	return &BaseTask{
		RunID:       primitive.NewObjectID().Hex(),
		WorkspaceID: workspaceID,
		ProcessorID: processorID,
		UploadID:    uploadID,
		TmpDir:      tmpDir,
		Input:       map[string]interface{}{"inputObjId": inputObjId},
		TmpDirLock:  &sync.Mutex{},
	}
}

func (t *BaseTask) MakeTask(data *BaseTask) {
	t.RunID = data.RunID
	t.TaskID = data.TaskID
	t.ProcessorID = data.ProcessorID
	t.WorkspaceID = data.WorkspaceID
	t.UploadID = data.UploadID
	t.TmpDir = data.TmpDir
	t.TmpDirLock = data.TmpDirLock
	t.TaskParams = data.TaskParams
	t.Input = data.Input
	t.Output = data.Output
}

func (t *BaseTask) GetTask() *BaseTask {
	return t
}
