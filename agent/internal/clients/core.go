package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/phuslu/log"
	"github.com/uploadpilot/agent/internal/dto"
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

func (c *CoreServiceClient) LogUploadRequest(ctx context.Context, originalReq *http.Request, workspaceID string, data map[string]string) (int, error) {
	logURI := fmt.Sprintf("/workspaces/%s/uploads/log", workspaceID)

	log := &dto.UploadRequestLog{
		WorkspaceID: workspaceID,
		Timestamp:   time.Now(),
		Data:        data,
	}

	body, err := json.Marshal(log)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resp, err := c.makeRequest(ctx, originalReq, "POST", logURI, bytes.NewReader(body))
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, fmt.Errorf("%s", responseBody)
	}

	return resp.StatusCode, nil
}

func (c *CoreServiceClient) GetUploaderConfig(ctx context.Context, originalReq *http.Request, workspaceID string) (*dto.WorkspaceConfig, int, error) {
	configURI := fmt.Sprintf("/workspaces/%s/config", workspaceID)
	resp, err := c.makeRequest(ctx, originalReq, "GET", configURI, nil)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf("%s", body)
	}

	var config dto.WorkspaceConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &config, resp.StatusCode, nil
}

func (c *CoreServiceClient) CreateNewUpload(ctx context.Context, originalReq *http.Request, workspaceID string, upload *dto.Upload) (string, int, error) {
	uploadID, statusCode, err := c.createUpload(ctx, originalReq, workspaceID, upload)
	if err != nil {
		return "", statusCode, err
	}

	return uploadID, statusCode, nil
}

func (c *CoreServiceClient) FinishUpload(ctx context.Context, originalReq *http.Request, workspaceID string, uploadID string, upload *dto.Upload) (int, error) {
	return c.patchUpload(ctx, originalReq, workspaceID, uploadID, upload)
}

func (c *CoreServiceClient) patchUpload(ctx context.Context, originalReq *http.Request, workspaceID string, uploadID string, upload *dto.Upload) (int, error) {
	finishUploadURI := fmt.Sprintf("/workspaces/%s/uploads/%s/finish", workspaceID, uploadID)
	body, err := json.Marshal(upload)
	if err != nil {
		return http.StatusBadRequest, err
	}
	resp, err := c.makeRequest(ctx, originalReq, "POST", finishUploadURI, bytes.NewReader(body))
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, fmt.Errorf("%s", responseBody)
	}

	return resp.StatusCode, nil
}

func (c *CoreServiceClient) createUpload(ctx context.Context, originalReq *http.Request, workspaceID string, upload *dto.Upload) (string, int, error) {
	createUploadURI := fmt.Sprintf("/workspaces/%s/uploads", workspaceID)
	body, err := json.Marshal(upload)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	resp, err := c.makeRequest(ctx, originalReq, "POST", createUploadURI, bytes.NewReader(body))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	if resp.StatusCode != 200 {
		return "", resp.StatusCode, fmt.Errorf("%s", responseBody)
	}

	var uploadID string
	err = json.Unmarshal(responseBody, &uploadID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return uploadID, resp.StatusCode, nil
}

func (c *CoreServiceClient) makeRequest(ctx context.Context, originReq *http.Request, method, endpoint string,
	body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseUri, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	apiKey := originReq.Header.Get("X-Api-Key")
	log.Debug().Str("api_key", apiKey).Msg("api key found")
	if apiKey != "" {
		req.Header.Set("X-Api-Key", apiKey)
	}

	req.Header.Set("X-Tenant-Id", originReq.Header.Get("X-Tenant-Id"))

	// copy all cookies from originalReq
	for _, cookie := range originReq.Cookies() {
		req.AddCookie(cookie)
	}

	return c.httpClient.Do(req)
}
