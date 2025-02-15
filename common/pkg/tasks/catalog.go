package tasks

import (
	commonutils "github.com/uploadpilot/uploadpilot/go-core/common/utils"
)

var tasks = map[string]interface{}{
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

func ValidateTaskData(key string, data map[string]interface{}) error {
	var newTask interface{}

	if err := commonutils.MapStructAndValidate(data, newTask); err != nil {
		return err
	}

	return nil
}
