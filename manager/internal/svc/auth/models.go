package auth

type WorkspacePerm string

const (
	CanRead   WorkspacePerm = "read_ws"
	CanManage WorkspacePerm = "manage_ws"
	CanUpload WorkspacePerm = "upload_ws"
)

type AccountPerm string

const (
	CanReadAcc   AccountPerm = "read_acc"
	CanManageAcc AccountPerm = "manage_acc"
)
