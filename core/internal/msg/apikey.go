package msg

// API request validation
const (
	ErrInvalidRequestBody     = "invalid request body. error: %s"
	ErrInvalidRequestParams   = "invalid request params. error: %s"
	ErrInvaliduestQueryParams = "invalid request query params. error: %s"
)

// API Key errors
const (
	ErrNoScopeInAPIKeyCreateRequest = "no scope in api key create request"
	ErrAPIKeyNotFoundInRequest      = "api key not found in request"
	ErrApiKeyCreateFailed           = "failed to create api key. please try again later"
	ErrInvalidAPIKey                = "invalid api key"
	ErrExpiredAPIKey                = "api key has expired"
	ErrRevokedAPIKey                = "api key has been revoked"
)
