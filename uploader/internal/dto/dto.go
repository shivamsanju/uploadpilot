package dto

import "time"

type WorkspaceConfig struct {
	MaxFileSize            *int64   `json:"maxFileSize"`
	MaxNumberOfFiles       int64    `json:"maxNumberOfFiles"`
	AllowedFileTypes       []string `json:"allowedFileTypes"`
	AllowedSources         []string `json:"allowedSources"`
	RequiredMetadataFields []string `json:"requiredMetadataFields"`
	AllowPauseAndResume    bool     `json:"allowPauseAndResume"`
	EnableImageEditing     bool     `json:"enableImageEditing"`
	UseCompression         bool     `json:"useCompression"`
	UseFaultTolerantMode   bool     `json:"useFaultTolerantMode"`
	ChunkSize              int64    `json:"chunkSize"`
	AllowedOrigins         []string `json:"-"`
}

type Upload struct {
	ID         string                 `json:"id"`
	FileName   string                 `json:"fileName"`
	FileType   string                 `json:"fileType"`
	Size       int64                  `json:"size"`
	Metadata   map[string]interface{} `json:"metadata"`
	Status     string                 `json:"status"`
	FinishedAt time.Time              `json:"finishedAt"`
	StartedAt  time.Time              `json:"startedAt"`
}

type UploadRequestLog struct {
	WorkspaceID string            `json:"workspaceId"`
	Timestamp   time.Time         `json:"timestamp"`
	Data        map[string]string `json:"data"`
}
