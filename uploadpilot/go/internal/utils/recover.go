package utils

import "github.com/uploadpilot/uploadpilot/internal/infra"

func Recover() {
	if r := recover(); r != nil {
		infra.Log.Errorf("recovered from panic: %s", r)
	}
}
