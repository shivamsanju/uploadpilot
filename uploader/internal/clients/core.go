package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uploadpilot/uploader/internal/dto"
)

type CoreServiceClient struct {
	httpClient *http.Client
	baseUri    string
}

func NewCoreServiceClient(baseUri string) *CoreServiceClient {
	return &CoreServiceClient{
		httpClient: &http.Client{},
		baseUri:    baseUri,
	}
}

func (c *CoreServiceClient) LogUploadRequest(ctx context.Context, workspaceID string, data map[string]string) error {
	logURI := fmt.Sprintf("/workspaces/%s/uploads/log", workspaceID)

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

func (c *CoreServiceClient) GetUploaderConfig(ctx context.Context, workspaceID string, headers http.Header) (*dto.WorkspaceConfig, error) {
	configURI := fmt.Sprintf("/workspaces/%s/config", workspaceID)
	resp, err := c.makeRequest(ctx, "GET", configURI, nil, headers)
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

	var config dto.WorkspaceConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *CoreServiceClient) CreateNewUpload(ctx context.Context, workspaceID string, upload *dto.Upload, headers http.Header) (string, error) {
	uploadID, err := c.createUpload(ctx, workspaceID, upload, headers)
	if err != nil {
		return "", err
	}

	return uploadID, nil
}

func (c *CoreServiceClient) FinishUpload(ctx context.Context, workspaceID string, uploadID string, upload *dto.Upload, headers http.Header) error {
	return c.patchUpload(ctx, workspaceID, uploadID, upload, headers)
}

func (c *CoreServiceClient) patchUpload(ctx context.Context, workspaceID string, uploadID string, upload *dto.Upload, headers http.Header) error {
	finishUploadURI := fmt.Sprintf("/workspaces/%s/uploads/%s/finish", workspaceID, uploadID)
	body, err := json.Marshal(upload)
	if err != nil {
		return err
	}
	resp, err := c.makeRequest(ctx, "POST", finishUploadURI, bytes.NewReader(body), headers)
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

func (c *CoreServiceClient) createUpload(ctx context.Context, workspaceID string, upload *dto.Upload, headers http.Header) (string, error) {
	createUploadURI := fmt.Sprintf("/workspaces/%s/uploads", workspaceID)
	body, err := json.Marshal(upload)
	if err != nil {
		return "", err
	}
	resp, err := c.makeRequest(ctx, "POST", createUploadURI, bytes.NewReader(body), headers)
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

func (c *CoreServiceClient) makeRequest(ctx context.Context, method, endpoint string,
	body io.Reader, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseUri, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	authHeader := headers.Get("Authorization")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	apiKey := headers.Get("X-Api-Key")
	if apiKey != "" {
		req.Header.Set("X-Api-Key", apiKey)
	}

	return c.httpClient.Do(req)
}
