package msg

var (
	DBConnectionSuccess = "successfully connected to postgres!"
	DBConnectionFailure = "failed to connect to postgres! err: %s"

	DBError             = "database error: %s"
	DBInvalidInput      = "invalid input: %s"
	DBInvalidID         = "invalid id: %s"
	DBErrRecordNotFound = "record not found"

	RedisConnectionSuccess = "successfully connected to redis!"
	RedisConnectionFailure = "failed to connect to redis! err: %s"

	RedisErrFailedToSet = "failed to set redis key: %s"
	RedisErrFailedToGet = "failed to get redis key: %s"
	RedisErrFailedToDel = "failed to delete redis key: %s"
)
