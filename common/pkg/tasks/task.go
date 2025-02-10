package tasks

import "github.com/uploadpilot/uploadpilot/common/pkg/models"

type Task[Data any] struct {
	Key         models.TaskKey `json:"key"`
	Label       string         `json:"label"`
	Description string         `json:"description"`
	Data        Data           `json:"message"`
}
