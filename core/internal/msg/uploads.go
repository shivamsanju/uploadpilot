package msg

const (
	ErrUploadSizeLessThanRequired           = "upload size less than required. size: %d, required: %d"
	ErrUploadSizeExceedAllowedLimit         = "upload size exceed allowed limit. size: %d, limit: %d"
	ErrUploadContentTypeNotAllowed          = "upload content type not allowed"
	ErrUploadMissingRequiredMetadataFields  = "upload missing required metadata fields: %s"
	ErrUploadURLValidityExceedsAllowedLimit = "requested upload url validity exceeds allowed limit. validity: %d, limit: %d"
	ErrUploadNotFinished                    = "upload not finished"
	ErrUploadAlreadyIsTerminalState         = "upload already is a terminal state"
)
