package uploader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type File struct {
	Name        string
	Data        []byte
	ContentType string
}

type Uploader struct {
	APIKey      string
	TenantID    string
	WorkspaceID string
	BaseURL     string
}

type UploaderOpts struct {
	BaseURL string
}

func NewUploader(tenantID, workspaceID, apiKey string, opts *UploaderOpts) *Uploader {
	if opts == nil {
		opts = &UploaderOpts{
			BaseURL: "http://localhost:8080",
		}
	}
	return &Uploader{
		TenantID:    tenantID,
		WorkspaceID: workspaceID,
		APIKey:      apiKey,
		BaseURL:     opts.BaseURL,
	}
}

func (u *Uploader) Upload(file *File, metadata map[string]interface{}) (bool, error) {
	if len(file.Data) == 0 {
		return false, errors.New("no file provided for upload")
	}

	uploadID, uploadURL, method, err := u.getPresignedUrl(file, metadata)
	if err != nil {
		return false, err
	}

	err = u.uploadToS3(file, uploadURL, method)
	if err != nil {
		return false, err
	}

	return u.completeUpload(uploadID)
}

func (u *Uploader) getPresignedUrl(file *File, metadata map[string]interface{}) (string, string, string, error) {
	url := fmt.Sprintf("%s/tenants/%s/workspaces/%s/uploads", u.BaseURL, u.TenantID, u.WorkspaceID)

	requestBody, err := json.Marshal(map[string]interface{}{
		"fileName":              file.Name,
		"contentType":           file.ContentType,
		"contentLength":         len(file.Data),
		"uploadUrlValiditySecs": 900,
		"metadata":              metadata,
	})
	if err != nil {
		return "", "", "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", u.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("failed to get upload URL: %s", respBody)
	}

	var result struct {
		UploadURL string `json:"uploadUrl"`
		Method    string `json:"method"`
		UploadID  string `json:"uploadId"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.UploadID, result.UploadURL, result.Method, nil
}

func (u *Uploader) uploadToS3(file *File, uploadURL, method string) error {
	req, err := http.NewRequest(method, uploadURL, bytes.NewReader(file.Data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", file.ContentType)
	req.Header.Set("If-None-Match", "*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload: %s", respBody)
	}

	return nil
}

func (u *Uploader) completeUpload(uploadID string) (bool, error) {
	if uploadID == "" {
		return false, errors.New("no uploadId provided")
	}

	url := fmt.Sprintf("%s/tenants/%s/workspaces/%s/uploads/%s/finish", u.BaseURL, u.TenantID, u.WorkspaceID, uploadID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("X-Api-Key", u.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to complete upload: %s", respBody)
	}

	return true, nil
}
