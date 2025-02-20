package templates

import "github.com/uploadpilot/uploadpilot/manager/internal/dto"

var ProcessorTemplates = []dto.ProcessorTemplate{
	{Key: "send_webhook", Label: "Send webhook", Description: "Sends a webhook to a target URL"},
	{Key: "send_email", Label: "Send email", Description: "Sends an email to a target email address"},
	{Key: "send_sms", Label: "Send SMS", Description: "Sends an SMS to a target phone number"},
	{Key: "send_slack_message", Label: "Send slack message", Description: "Sends a message to a Slack channel"},
	{Key: "extract_content_from_pdf", Label: "Extract content from a PDF", Description: "Extract xontent from pdf files"},
}
