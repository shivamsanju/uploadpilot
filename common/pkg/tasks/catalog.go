package tasks

import (
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
)

var tasks = map[models.TaskKey]interface{}{
	WebhookTask.Key:           WebhookTask,
	ExtractPDFContentTask.Key: ExtractPDFContentTask,
}

func GetAllTasks() []interface{} {
	var tasksList []interface{}
	for _, task := range tasks {
		tasksList = append(tasksList, task)
	}
	return tasksList
}

func ValidateTaskData(key models.TaskKey, data map[string]interface{}) error {
	var newTask interface{}

	if err := commonutils.MapStructAndValidate(data, newTask); err != nil {
		return err
	}

	return nil
}
