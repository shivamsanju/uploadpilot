package hooks

import (
	"fmt"
	"net/http"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

func AuthHook(hook tusd.HookEvent) error {
	headers := hook.HTTPRequest.Header
	token := headers.Get("Authorization")
	if len(token) == 0 {
		return fmt.Errorf("missing bearer token in header")
	}
	// Get auth hook - send a http request
	req, err := http.NewRequest("GET", "http://localhost:3001/api/auth", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		infra.Log.Errorf("auth hook failed: %s", err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		infra.Log.Errorf("auth hook failed: %s", resp.Status)
		return fmt.Errorf("auth hook failed: %s", resp.Status)
	}
	defer resp.Body.Close()
	return nil
}
