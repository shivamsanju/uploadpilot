package utils

import "github.com/uploadpilot/uploadpilot/common/pkg/infra"

func Recover() {
	if r := recover(); r != nil {
		infra.Log.Errorf("recovered from panic: %s", r)
	}
}
