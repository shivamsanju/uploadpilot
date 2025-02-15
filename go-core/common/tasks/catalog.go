package tasks

var TaskCatalog = map[string]*Task{
	"Webhook":           WebhookTask,
	"ExtractPDFContent": ExtractPDFContentTask,
}
