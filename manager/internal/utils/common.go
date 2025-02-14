package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/uploadpilot/uploadpilot/common/pkg/msg"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
)

func GetUserDetailsFromContext(ctx context.Context) (*dto.UserContext, error) {
	userID, ok1 := ctx.Value(dto.UserIDContextKey).(string)
	name, ok2 := ctx.Value(dto.NameContextKey).(string)
	email, ok3 := ctx.Value(dto.EmailContextKey).(string)

	if !ok1 || !ok2 || !ok3 {
		return nil, errors.New(msg.FailedToGetUserFromContext)
	}

	return &dto.UserContext{
		UserID: userID,
		Name:   name,
		Email:  email,
	}, nil
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

func GetSkipLimitSearchParams(r *http.Request) (skip int, limit int, search string, err error) {
	query := r.URL.Query()
	s := query.Get("skip")
	l := query.Get("limit")
	search = query.Get("search")

	skip, err = strconv.Atoi(s)
	if err != nil {
		return
	}
	limit, err = strconv.Atoi(l)
	if err != nil {
		return
	}
	return
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
