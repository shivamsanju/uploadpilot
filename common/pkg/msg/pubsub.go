package msg

const (
	RetryingFailedMessages = "retrying [%d] failed messages"
	ErrReadingFromStream   = "error reading from stream: %s"
	ErrClaimingMessage     = "error claiming message: %s"

	ErrInvalidMessageFormat = "message [%s] has an invalid message format"
	ErrDecodingMessage      = "error decoding message [%s]: %s"
	ErrorInHandler          = "error while handling message [%s]: %s"
)
