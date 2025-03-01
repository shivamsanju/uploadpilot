package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type ErrorResponse struct {
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

func HandleHttpError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	reqID := middleware.GetReqID(r.Context())
	log.Error().Err(err).Str("request_id", reqID).Int("status_code", statusCode).Msg("request failed")
	render.Status(r, statusCode)
	render.JSON(w, r, &ErrorResponse{
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

func VerifyWorkspaceID(workspaceID string) error {
	if workspaceID == "" {
		return errors.New("workspaceId is required")
	}

	if _, err := uuid.Parse(workspaceID); err != nil {
		return errors.New("invalid workspaceId")
	}

	return nil
}
