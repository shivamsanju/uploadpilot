package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/uploader/internal/dto"
)

func HandleHttpError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	reqID := middleware.GetReqID(r.Context())
	infra.Log.Errorf("request with id [%s] failed: %s", reqID, err.Error())
	render.Status(r, statusCode)
	infra.Log.Infof("STATUS: %d", statusCode)
	render.JSON(w, r, &dto.ErrorResponse{
		RequestID: reqID,
		Message:   err.Error(),
	})
}

func GetStatusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return fmt.Sprintf("%d OK", status)
	case status >= 300 && status < 400:
		return fmt.Sprintf("%d Redirect", status)
	case status >= 400 && status < 500:
		return fmt.Sprintf("%d Client Error", status)
	case status >= 500:
		return fmt.Sprintf("%d Server Error", status)
	default:
		return fmt.Sprintf("%d Unknown", status)
	}
}
