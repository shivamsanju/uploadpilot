package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
)

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
