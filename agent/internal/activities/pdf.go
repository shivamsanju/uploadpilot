package activities

import (
	"context"
	"fmt"
	"math/rand"

	"go.temporal.io/sdk/activity"
)

type PdfActivities struct {
}

func (a *PdfActivities) ExtractContentFromPDF(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	randomBool := rand.Intn(3) == 0 // Generates 0 or 1, maps to true/false
	fmt.Println(randomBool)
	if randomBool {
		fmt.Printf("Extracting content success from PDF %s with input %v \n", name, argsStr)
		return "Extracted content successfully", nil
	} else {
		fmt.Printf("Extracting content failed from PDF %s with input %v \n", name, argsStr)
		return "Extraction failed", fmt.Errorf("extraction failed")
	}
}
