package msg

// Common errors
const (
	ErrUnexpected = "oops! something went wrong. please try again"
	ErrUnknown    = "there was an issue processing your request. please try again later"
)

// User
const (
	ErrUserInfoNotFoundInRequest = "user info not found in request"
	ErrInvalidUserInfoInRequest  = "invalid user info in request"
	ErrTenantIDNotFoundInRequest = "tenant id not found in request"
	ErrInvalidTenantIDInRequest  = "invalid tenant id in request"
	ErrEmailNotFoundInRequest    = "email not found in request"
)

const (
	//Error

	// Structs
	ErrFailedToMarshal = "failed to marshal value: %s"

	// API
	ValidationErr              = "validation error: %s"
	FailedToGetUserFromContext = "failed to get user from request"

	// Workspace
	WorkspaceNotFound            = "workspace %s not found"
	UserAlreadyExistsInWorkspace = "user %s already exists in workspace"
	UnknownRole                  = "unknown role: %s"
	OwnerCannotBeRemoved         = "owner cannot be removed"
	OwnerRoleCannotBeChanged     = "owner role cannot be changed"
	UserNotFound                 = "user %s not found"

	// Upload
	UploadNotFound              = "upload %s not found"
	InvalidWorkspaceIDInHeaders = "invalid workspace id in headers: %s"

	// Processors
	ProcTaskFailed = "task %s failed for workspaceID %s and processorID %s and uploadID %s. error: %s"

	ProcessingFailed    = "upload processing failed for processorID %s"
	ProcessingComplete  = "upload processing complete for processorID %s"
	ProcessingCancelled = "upload processing cancelled for processorID %s"

	// Upload Validations
	UploadWithinSizeLimits = "upload size: %d is within allowed range: %d to %s"
	MaxFileSizeExceeded    = "max file size exceeded: %d > %d"
	MinFileSizeNotMet      = "min file size not met: %d < %d"

	FileTypeValidationSkipped = "file type validation skipped as all file types are allowed"
	FileTypeValidationPassed  = "file type validation passed. [%s] is among the allowed file types"
	FileTypeValidationFailed  = "file type validation failed. [%s] is not among the allowed file types"

	UploadAuthenticationSkipped = "upload authentication skipped as auth endpoint is not set"
	UploadAuthenticationPassed  = "upload authentication passed"
	UploadAuthenticationFailed  = "upload authentication failed. reason: %s"

	// Upload Status
	FailedToCreateUpload   = "failed to create upload. reason: %s"
	UploadStarted          = "upload started"
	UploadFailed           = "upload failed. reason: %s"
	UploadComplete         = "upload completed"
	ErrSubscriptionExpired = "your subscription has expired. please renew your subscription to continue using uploadpilot"
)
