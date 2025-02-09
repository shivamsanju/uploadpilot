package proc

import (
	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/proc/tasks"
)

var taskRegistry = make(map[models.TaskKey]TaskBlock)

func InitProcRegistry() {
	for _, block := range ProcTaskBlocks {
		taskRegistry[block.Key] = block
	}
}

func GetTaskFromRegistry(key models.TaskKey) (tasks.Task, bool) {
	if block, exists := taskRegistry[key]; exists {
		return block.TaskBuilder(), true
	}
	return nil, false
}
