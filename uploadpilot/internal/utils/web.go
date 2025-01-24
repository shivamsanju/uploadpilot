package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
)

func GetSkipLimitSearchParams(r *http.Request) (skip int64, limit int64, search string, err error) {
	query := r.URL.Query()
	s := query.Get("skip")
	l := query.Get("limit")
	search = query.Get("search")

	skip, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	limit, err = strconv.ParseInt(l, 10, 64)
	if err != nil {
		return
	}
	return
}

func GetUserDetailsFromContext(ctx context.Context) (*dto.ApiUser, error) {
	userID, ok1 := ctx.Value("userId").(string)
	name, ok2 := ctx.Value("name").(string)
	email, ok3 := ctx.Value("email").(string)

	if !ok1 || !ok2 || !ok3 {
		return nil, errors.New(msg.FailedToGetUserFromContext)
	}

	return &dto.ApiUser{
		UserID: userID,
		Name:   name,
		Email:  email,
	}, nil
}

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

func ExtractKeyValuePairs(search string) (map[string]string, error) {
	params := make(map[string]string)
	search = strings.TrimSpace(search)
	search = strings.ReplaceAll(search, "{", "")
	search = strings.ReplaceAll(search, "}", "")
	pairs := strings.Split(search, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid search format: %s", pair)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		params[key] = value
	}

	return params, nil
}
