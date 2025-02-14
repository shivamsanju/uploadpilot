package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
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
