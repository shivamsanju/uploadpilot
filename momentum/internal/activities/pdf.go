package activities

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
)

type PdfActivities struct {
}

func (a *PdfActivities) ExtractContentFromPDF(ctx context.Context, input []string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Extracting content from PDF %s with input %v \n", name, input)
	return "Extracted content from PDF_" + name, nil
}
