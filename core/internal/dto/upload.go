package dto

type CreateUploadRequest struct {
	FileName              string                 `json:"fileName" validate:"required"`
	ContentType           string                 `json:"contentType" validate:"required"`
	ContentLength         int64                  `json:"contentLength" validate:"required"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
	UploadURLValiditySecs int64                  `json:"uploadUrlValiditySecs" validate:"required"`
}

type CreateUploadResponse struct {
	UploadID      string              `json:"uploadId" validate:"required"`
	UploadURL     string              `json:"uploadUrl" validate:"required"`
	Method        string              `json:"method" validate:"required"`
	SignedHeaders map[string][]string `json:"signedHeaders" validate:"required"`
}
