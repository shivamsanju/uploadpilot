package sdk

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/eventials/go-tus"
)

type UploadOptions struct {
	Metadata map[string]string
	Headers  map[string]string
}

type Uploader interface {
	UploadFile(path string, options *UploadOptions) error
	Upload(data []byte, options *UploadOptions) error
}

type uploader struct {
	workspaceID string
	baseURI     string
}

func NewUploader(workspaceID string, baseURI string) (Uploader, error) {
	if workspaceID == "" || baseURI == "" {
		return nil, errors.New("invalid parameters: workspaceID, baseURI are required")
	}
	return &uploader{
		workspaceID: workspaceID,
		baseURI:     baseURI,
	}, nil
}

// UploadFile uploads a file using the tus protocol with metadata and headers.
func (t *uploader) UploadFile(path string, options *UploadOptions) error {
	if path == "" || options == nil {
		return errors.New("invalid parameters:  path, and options are required")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err.Error())
	}
	defer f.Close()

	client, err := createTusClient(getBaseURI(t.baseURI, t.workspaceID), options.Headers)
	if err != nil {
		return err
	}

	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		return err
	}

	if err := setUploadMetadata(upload, options.Metadata); err != nil {
		return err
	}

	uploader, err := client.CreateUpload(upload)
	if err != nil {
		return err
	}

	err = uploader.Upload()
	if err != nil {
		return err
	}

	return nil
}

// Upload uploads byte data using the tus protocol with metadata and headers.
func (t *uploader) Upload(data []byte, options *UploadOptions) error {
	if data == nil || options == nil {
		return errors.New("invalid parameters:  data, and options are required")
	}

	filename, err := createTempFile(data)
	if err != nil {
		return err
	}
	defer os.Remove(filename)

	return t.UploadFile(filename, options)
}

// createTusClient initializes a tus client with custom headers.
func createTusClient(tusServerURL string, headers map[string]string) (*tus.Client, error) {
	options := tus.Config{
		ChunkSize:           tus.DefaultConfig().ChunkSize,
		Header:              make(http.Header),
		Resume:              false,
		OverridePatchMethod: false,
		Store:               nil,
		HttpClient:          nil,
	}
	for key, value := range headers {
		options.Header.Set(key, value)
	}
	client, err := tus.NewClient(tusServerURL, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to create tus client: %w", err)
	}
	return client, nil
}

// createTempFile creates a temporary file and writes data to it.
func createTempFile(data []byte) (string, error) {
	contentType := http.DetectContentType(data)
	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extensions) == 0 {
		return "", fmt.Errorf("failed to detect content type: %w", err)
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extensions[0])
	f, err := os.CreateTemp("", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data to file: %w", err)
	}

	return f.Name(), nil
}

// setUploadMetadata attaches metadata to the upload object.
func setUploadMetadata(upload *tus.Upload, metadata map[string]string) error {
	_, ok := metadata["filename"]
	if !ok {
		return errors.New("filename not found in metadata")
	}

	_, ok = metadata["filetype"]
	if !ok {
		return errors.New("filetype not found in metadata")
	}

	if upload.Metadata == nil {
		upload.Metadata = make(map[string]string)
	}
	for key, value := range metadata {
		upload.Metadata[key] = value
	}

	return nil

}

// getBaseURI returns a valid base URI, defaulting to localhost if empty.
func getBaseURI(baseURI string, workspaceID string) string {
	if baseURI == "" {
		return "http://localhost:8081/upload/" + workspaceID
	}
	return baseURI + "/upload/" + workspaceID
}
