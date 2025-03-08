package workflow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/phuslu/log"
)

type Executor struct {
	lambdaClient *lambda.Client
}

func NewExecutor(lambdaClient *lambda.Client) *Executor {
	return &Executor{
		lambdaClient: lambdaClient,
	}
}

func (e *Executor) ExecuteLambdaContainerActivity(ctx context.Context, functionName, marshaledPayload string) ([]byte, error) {

	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      []byte(marshaledPayload),
		LogType:      types.LogTypeTail,
	}

	log.Info().Str("functionName", functionName).Str("payload", marshaledPayload).Msg("invoking lambda")
	op, err := e.lambdaClient.Invoke(context.TODO(), input)
	if err != nil {
		log.Error().Err(err).Msg("failed to invoke lambda")
		return nil, fmt.Errorf("failed to run activity: %w", err)
	}

	if op.FunctionError != nil {
		log.Error().Str("error", *op.FunctionError).Msg("lambda error")
		return nil, fmt.Errorf("error: %s", string(op.Payload))
	}

	var output map[string]interface{}
	if err := json.Unmarshal(op.Payload, &output); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal lambda output")
		return nil, fmt.Errorf("failed to unmarshal lambda output: %w", err)
	}

	success, ok := output["success"]
	if !ok {
		log.Error().Str("output", string(op.Payload)).Msg("lambda output")
		return nil, fmt.Errorf("failed to unmarshal lambda output: success field not found")
	}

	s, ok := success.(bool)
	if !ok || !s {
		err := output["error"].(string)
		log.Error().Str("error", err).Msg("lambda execution failed")
		return nil, fmt.Errorf("error: %s", err)
	}

	if op.FunctionError != nil {
		log.Error().Str("error", *op.FunctionError).Msg("lambda error")
		return nil, fmt.Errorf("error: %s", string(op.Payload))
	}

	log.Info().Str("output", string(op.Payload)).Msg("lambda output")
	logs, _ := DecodeBase64(*op.LogResult)
	log.Info().Str("logs", logs).Msg("lambda logs")

	return op.Payload, nil
}

// DecodeBase64 takes a Base64-encoded string and returns the decoded string
func DecodeBase64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
