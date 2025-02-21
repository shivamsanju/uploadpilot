package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/phuslu/log"
	"github.com/uploadpilot/uploader/internal/dto"
)

type CoreServiceClient struct {
	httpClient *http.Client
	baseUri    string
	apiKey     string
}

func NewCoreServiceClient(baseUri string, apiKey string) *CoreServiceClient {
	return &CoreServiceClient{
		httpClient: &http.Client{},
		baseUri:    baseUri,
		apiKey:     apiKey,
	}
}

func (c *CoreServiceClient) LogUploadRequest(ctx context.Context, workspaceID string, data map[string]string) error {
	logURI := fmt.Sprintf("/workspaces/%s/log-upload", workspaceID)

	log := &dto.UploadRequestLog{
		WorkspaceID: workspaceID,
		Timestamp:   time.Now(),
		Data:        data,
	}

	body, err := json.Marshal(log)
	if err != nil {
		return err
	}

	resp, err := c.makeRequest(ctx, "POST", logURI, bytes.NewReader(body), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", responseBody)
	}

	return nil
}

func (c *CoreServiceClient) GetUploaderConfig(ctx context.Context, workspaceID string) (*dto.UploaderConfig, error) {
	configURI := fmt.Sprintf("/workspaces/%s/config", workspaceID)
	resp, err := c.makeRequest(ctx, "GET", configURI, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	var config dto.UploaderConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *CoreServiceClient) CreateNewUpload(ctx context.Context, workspaceID string, upload *dto.Upload) (string, error) {
	sub, err := c.verifySubscription(ctx, workspaceID)
	if err != nil {
		log.Error().Err(err).Str("workspace_id", workspaceID).Msg("failed to get subscription information")
		return "", errors.New("failed to get subscription information")
	}
	if sub == nil || !sub.Active {
		return "", errors.New("no active subscription found")
	}

	uploadID, err := c.createUpload(ctx, workspaceID, upload)
	if err != nil {
		return "", err
	}

	return uploadID, nil
}

func (c *CoreServiceClient) FinishUpload(ctx context.Context, workspaceID string, uploadID string, upload *dto.Upload) error {
	return c.patchUpload(ctx, workspaceID, uploadID, upload)
}

func (c *CoreServiceClient) patchUpload(ctx context.Context, workspaceID string, uploadID string, upload *dto.Upload) error {
	finishUploadURI := fmt.Sprintf("/workspaces/%s/uploads/%s/finish", workspaceID, uploadID)
	body, err := json.Marshal(upload)
	if err != nil {
		return err
	}
	resp, err := c.makeRequest(ctx, "POST", finishUploadURI, bytes.NewReader(body), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", responseBody)
	}

	return nil
}

func (c *CoreServiceClient) createUpload(ctx context.Context, workspaceID string, upload *dto.Upload) (string, error) {
	createUploadURI := fmt.Sprintf("/workspaces/%s/uploads", workspaceID)
	body, err := json.Marshal(upload)
	if err != nil {
		return "", err
	}
	resp, err := c.makeRequest(ctx, "POST", createUploadURI, bytes.NewReader(body), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s", responseBody)
	}

	var uploadID string
	err = json.Unmarshal(responseBody, &uploadID)
	if err != nil {
		return "", err
	}

	return uploadID, nil
}

func (c *CoreServiceClient) verifySubscription(ctx context.Context, workspaceID string) (*dto.Subscription, error) {
	verifySubscriptionURI := fmt.Sprintf("/workspaces/%s/subscription", workspaceID)
	resp, err := c.makeRequest(ctx, "GET", verifySubscriptionURI, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	var subscription dto.Subscription
	err = json.Unmarshal(body, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (c *CoreServiceClient) makeRequest(ctx context.Context, method, endpoint string,
	body io.Reader, headers map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseUri, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "uploadpilot-uploader")
	req.Header.Set("X-Api-Key", c.apiKey)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.httpClient.Do(req)
}
