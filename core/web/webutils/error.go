package webutils

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/phuslu/log"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
)

func HandleHttpError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	if err.Error() == msg.ErrAccessDenied {
		statusCode = http.StatusForbidden
	}
	render.Status(r, statusCode)

	reqID := middleware.GetReqID(r.Context())
	log.Error().Msgf("request with id [%s] failed: %s", reqID, err.Error())

	render.JSON(w, r, &dto.ErrorResponse{
		RequestID: reqID,
		Message:   err.Error(),
	})
}
