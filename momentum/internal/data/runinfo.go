package data

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

type ActivityInfo struct {
	ActivityKey     string     `json:"activityKey"`
	DataContainerID string     `json:"dataContainerId"`
	NumFiles        int        `json:"numFiles"`
	NumBytes        int64      `json:"numBytes"`
	Files           []FileInfo `json:"fileNames"`
}

type RunInfo struct {
	WorkflowID      string                  `json:"workflowId"`
	RunID           string                  `json:"runId"`
	UploadID        string                  `json:"uploadId"`
	WorkspaceID     string                  `json:"workspaceId"`
	ProcessorID     string                  `json:"processorId"`
	ActivityInfoMap map[string]ActivityInfo `json:"activityInfoMap"`
}
