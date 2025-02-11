package main

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/uploadpilot/uploadpilot/common/pkg/cache"
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/kms"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := infra.Init(&infra.S3Config{
		AccessKey: config.S3AccessKey,
		SecretKey: config.S3SecretKey,
		Region:    config.S3Region,
	}); err != nil {
		panic(err)
	}

	err := db.Init(config.PostgresURI)

	if err != nil {
		panic(err)
	}

	err = kms.Init(config.EncryptionKey)

	if err != nil {
		panic(err)
	}

	// Initialize cache.
	if err := cache.Init(&config.RedisAddr, &config.RedisPassword, &config.RedisUsername, config.RedisTLS); err != nil {
		panic(err)
	}

	// Add a map to processor
	proc := models.Processor{
		WorkspaceID: "856a3bda-5e6f-4353-8da6-864f58136c6e",
		ID:          "123e4567-e89b-12d3-a456-426655440000",
		Name:        "Test Processor",
		Triggers:    []string{"testTrigger"},
		Tasks:       []models.Task{},
		Statement: map[string]interface{}{
			"sequence": map[string]interface{}{
				"elements": []map[string]interface{}{
					{
						"activity": map[string]interface{}{
							"id":        "123e4567-e89b-12d3-a456-426655440000",
							"label":     "Sample Activity",
							"key":       "SampleActivity1",
							"keys":      []string{"SampleActivity2", "SampleActivity3"},
							"arguments": []string{"arg1"},
							"result":    "result1",
						},
					},
					{
						"activity": map[string]interface{}{
							"id":        "123e4567-e89b-12d3-a456-426655440001",
							"label":     "Sample Activity 2",
							"key":       "SampleActivity2",
							"arguments": []string{"result1"},
							"result":    "result2",
						},
					},
					{
						"activity": map[string]interface{}{
							"id":        "123e4567-e89b-12d3-a456-426655440002",
							"label":     "Sample Activity 3",
							"key":       "SampleActivity3",
							"arguments": []string{"arg2", "result2"},
							"result":    "result3",
						},
					},
				},
			},
		},
		Variables: map[string]interface{}{
			"testVariable":  "testValue",
			"testVariable2": "testValue2",
		},
		Enabled: true,
	}

	pr := db.NewProcessorRepo()
	if err := pr.Create(context.Background(), &proc).Error; err != nil {
		panic(err)
	}

	// get the task
	procs, err := pr.Get(context.Background(), "123e4567-e89b-12d3-a456-426655440000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", procs)

	var statement models.Statement

	if err := mapstructure.Decode(procs.Statement, &statement); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", statement)

	fmt.Print(statement.Sequence.Elements[0].Activity.Key)

}
