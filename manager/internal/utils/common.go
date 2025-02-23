package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/msg"
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
func GetPaginatedQueryParams(query *dto.PaginatedQuery) (*models.PaginationParams, error) {
	if query == nil {
		return &models.PaginationParams{
			Offset: 0,
			Limit:  0,
			Sort:   models.Sort{},
			Filter: []models.Filter{},
			Search: "",
		}, nil
	}

	var offset int = 0
	if query.Offset != "" {
		var err error
		offset, err = strconv.Atoi(query.Offset)
		if err != nil {
			return nil, err
		}
	}

	var limit int = 0
	if query.Limit != "" {
		var err error
		limit, err = strconv.Atoi(query.Limit)
		if err != nil {
			return nil, err
		}
	}

	var search string
	if query.Search != "" {
		search = strings.TrimSpace(query.Search)
	}

	var filter []models.Filter
	if query.Filter != "" {
		filterKeyValues := strings.Split(query.Filter, ";")
		for _, kv := range filterKeyValues {
			kvals := strings.Split(kv, ":")
			if len(kvals) != 2 {
				return nil, fmt.Errorf("invalid filter format: %s", kv)
			}
			values := strings.Split(kvals[1], ",")
			if len(values) == 0 || values[0] == "" {
				return nil, fmt.Errorf("invalid filter format: %s", kv)
			}
			filter = append(filter, models.Filter{
				Field: kvals[0],
				Value: values,
			})
		}
	}

	var sort models.Sort
	if query.Sort != "" {
		sortValues := strings.Split(query.Sort, ":")
		if len(sortValues) != 2 || strings.TrimSpace(sortValues[0]) == "" {
			return nil, fmt.Errorf("invalid sort format: %s", query.Sort)
		}
		order := strings.TrimSpace(sortValues[1])
		if order != "asc" && order != "desc" {
			return nil, fmt.Errorf("invalid sort order: %s", query.Sort)
		}
		sort = models.Sort{
			Field: strings.TrimSpace(sortValues[0]),
			Order: order,
		}
	}

	return &models.PaginationParams{
		Offset:              offset,
		Limit:               limit,
		Search:              search,
		CaseSensitiveSearch: query.CaseSensitiveSearch == "true",
		Filter:              filter,
		Sort:                sort,
	}, nil
}
