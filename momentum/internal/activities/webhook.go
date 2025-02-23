package activities

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
)

type WebhookActivities struct {
}

func (a *WebhookActivities) SendWebhook(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Sending webhook %s with input %v \n", name, argsStr)
	return "Sending webhook to " + inputActivityKey, nil
}
