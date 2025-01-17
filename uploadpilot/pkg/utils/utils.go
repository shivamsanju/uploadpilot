package utils

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"github.com/uploadpilot/uploadpilot/pkg/web/models"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleHttpError(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	reqID := middleware.GetReqID(r.Context())
	g.Log.Errorf("request with id [%s] failed: %s", reqID, err.Error())
	render.Status(r, statusCode)
	render.JSON(w, r, &models.ErrorResponse{RequestID: reqID, Message: err.Error()})
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

func FilterNonEmptyBsonFields(data bson.M) bson.M {
	filtered := bson.M{}
	for key, value := range data {
		if value != nil && value != "" {
			filtered[key] = value
		}
	}
	return filtered
}
