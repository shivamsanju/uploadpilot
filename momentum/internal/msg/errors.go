package msg

const (
	ErrUnknown                    = "unknown_error"
	ErrTmpDirCreationFailed       = "err_tmp_dir_failed"
	ErrInputDownloadFailed        = "err_input_download_failed"
	ErrOutputSaveFailed           = "err_output_save_failed"
	ErrRunInfoFileNotFound        = "err_run_info_file_not_found"
	ErrRunInfoFileReadFailed      = "err_run_info_file_read_failed"
	ErrRunInfoFileUnmarshalFailed = "err_run_info_file_format_error"
	ErrRunInfoFileMarshalFailed   = "err_run_info_file_marshal_failed"
	ErrInvalidWorkflowMetadata    = "err_invalid_workflow_metadata"
	ErrRunInfoFileSaveFailed      = "err_run_info_file_save_failed"
	ErrTaskInfoNotFound           = "err_task_info_not_found"
	ErrInvalidActivityArguments   = "err_invalid_activity_arguments"
	ErrRunInfoSaveFailed          = "err_run_info_save_failed"
)
