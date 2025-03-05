package webutils

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/phuslu/log"
	"github.com/uploadpilot/core/internal/dto"
)

func HandleHttpError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	reqID := middleware.GetReqID(r.Context())
	log.Error().Msgf("request with id [%s] failed: %s", reqID, err.Error())
	render.Status(r, statusCode)
	log.Info().Msgf("STATUS: %d", statusCode)
	render.JSON(w, r, &dto.ErrorResponse{
		RequestID: reqID,
		Message:   err.Error(),
	})
}
