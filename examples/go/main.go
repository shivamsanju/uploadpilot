package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Client struct {
	BaseURL     string
	TenantID    string
	WorkspaceID string
	HTTPClient  *http.Client
	Headers     map[string]string
}

type CreateUploadRequest struct {
	FileName              string                 `json:"fileName"`
	ContentType           string                 `json:"contentType"`
	ContentLength         int64                  `json:"contentLength"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
	UploadURLValiditySecs int64                  `json:"uploadUrlValiditySecs"`
}

type CreateUploadResponse struct {
	UploadID      string              `json:"uploadId"`
	UploadURL     string              `json:"uploadUrl"`
	Method        string              `json:"method"`
	SignedHeaders map[string][]string `json:"signedHeaders"`
}

func NewClient(baseURL, tenantID, workspaceID string, headers map[string]string) *Client {
	return &Client{
		BaseURL:     baseURL,
		TenantID:    tenantID,
		WorkspaceID: workspaceID,
		HTTPClient:  &http.Client{},
		Headers:     headers,
	}
}

func (c *Client) requestUploadURL(fileName, contentType string, fileSize int64, validitySecs int64) (*CreateUploadResponse, error) {
	url := fmt.Sprintf("%s/tenants/%s/workspaces/%s/uploads", c.BaseURL, c.TenantID, c.WorkspaceID)
	requestBody := CreateUploadRequest{
		FileName:              fileName,
		ContentType:           contentType,
		ContentLength:         fileSize,
		UploadURLValiditySecs: validitySecs,
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", "up-Ls189TgdKIer5eBp20250314062549")
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyMap := make(map[string]interface{})
		if err := json.NewDecoder(resp.Body).Decode(&bodyMap); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create upload: %s. Reason: %s", resp.Status, bodyMap)
	}

	var uploadResp CreateUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	return &uploadResp, nil
}

func (c *Client) uploadFile(filePath string, uploadResp *CreateUploadResponse) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	log.Printf("Uploading %s (%d bytes)", fileInfo.Name(), fileInfo.Size())

	req, err := http.NewRequest(uploadResp.Method, uploadResp.UploadURL, file)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "image/png")
	req.Header.Set("If-None-Match", "*")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		bodyStr := string(bodyBytes)
		return fmt.Errorf("failed to upload file: %s. Reason: %s", resp.Status, bodyStr)
	}

	return nil
}

func (c *Client) completeUpload(uploadID string) error {
	url := fmt.Sprintf("%s/tenants/%s/workspaces/%s/uploads/%s/finish", c.BaseURL, c.TenantID, c.WorkspaceID, uploadID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "image/png")
	req.Header.Set("X-Api-Key", "up-Ls189TgdKIer5eBp20250314062549")
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to complete upload")
	}

	return nil
}

func (c *Client) UploadFile(filePath string, validitySecs int64) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	log.Printf("Uploading %s (%d bytes)", fileInfo.Name(), fileInfo.Size())

	uploadResp, err := c.requestUploadURL(fileInfo.Name(), "image/png", fileInfo.Size(), validitySecs)
	if err != nil {
		return err
	}

	if err := c.uploadFile(filePath, uploadResp); err != nil {
		return err
	}

	return c.completeUpload(uploadResp.UploadID)
}

func main() {
	cl := NewClient(
		"http://localhost:8080",
		"68b167fb-8b74-4716-999d-421c90ea9f8a",
		"f2f3ccc0-be46-422b-a77e-f57da6ab2262",
		map[string]string{},
	)
	if err := cl.UploadFile("x/error-dark.png", 60); err != nil {
		log.Fatal(err)
	}

	fmt.Println("file uploaded successfully")
}
