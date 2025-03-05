package rbac

type AppRole string

const (
	Admin    AppRole = "admin"
	Reader   AppRole = "reader"
	Uploader AppRole = "uploader"
)
