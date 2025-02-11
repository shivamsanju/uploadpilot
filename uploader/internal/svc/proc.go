package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
	"github.com/uploadpilot/uploadpilot/uploader/internal/dto"
)

func (s *Service) GetProcessors(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	processors, err := s.procRepo.GetAll(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return processors, nil
}

func (s *Service) BuildAndTriggerTasks(ctx context.Context, workspaceID string) error {
	processors, err := s.GetProcessors(ctx, workspaceID)
	if err != nil {
		return err
	}

	for _, processor := range processors {
		if err := s.TriggerTask(ctx, processor.Workflow); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) TriggerTask(ctx context.Context, yaml string) error {
	body, err := json.Marshal(&dto.TriggerworkflowReq{
		Workflow: yaml,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.MomentumEndpoint+"/trigger", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Secret", config.MomentumSecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("momentum responded with status code %d", resp.StatusCode)
	}

	return nil
}
