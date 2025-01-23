package messages

var (
	// JWT
	JWTSecretKeyNotSet = "secret key has not been set. call auth.Init() first"
	InvalidToken       = "invalid jwt token"
	TokenExpired       = "jwt token is expired"

	// API
	ValidationErr              = "validation error: %s"
	FailedToGetUserFromContext = "failed to get user from request"

	// Workspace
	WorkspaceNotFound            = "workspace %s not found"
	UserAlreadyExistsInWorkspace = "user %s already exists in workspace"
	UnknownRole                  = "unknown role: %s"
	LastOwnerCannotBeRemoved     = "atleast one owner must remain in the workspace"
	LastOwnerRoleCannotBeChanged = "last owner role cannot be changed"
	UserNotFound                 = "user %s not found"

	// DB
	InvalidObjectID = "invalid object id: %s"

	// Upload
	UploadNotFound              = "upload %s not found"
	MaxFileSizeExceeded         = "max file size exceeded: %d > %d"
	MinFileSizeNotMet           = "min file size not met: %d < %d"
	MaxNumFilesExceeded         = "max number of files exceeded: %d > %d"
	MinNumFilesNotMet           = "min number of files not met: %d < %d"
	InvalidWorkspaceIDInHeaders = "invalid workspace id in headers: %s"
	FileTypeNotAllowed          = "file type not allowed: %s"
	UploadAuthenticationFailed  = "authentication failed for upload request: %s"
)
