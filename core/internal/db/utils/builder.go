package dbutils

import (
	"fmt"
	"slices"

	"github.com/uploadpilot/core/internal/db/models"
	"gorm.io/gorm"
)

type PaginationQueryInput struct {
	PaginationParams    *models.PaginationParams
	AllowedSearchFields []string
	AllowedFilterFields []string
}

func BuildPaginationQuery(
	query *gorm.DB,
	input *PaginationQueryInput,
) (*gorm.DB, int64, bool, error) {
	var totalRecords int64
	var sortApplied bool

	if input.PaginationParams.Search != "" && len(input.AllowedSearchFields) > 0 {
		searchClause := "%" + input.PaginationParams.Search + "%"
		matchType := "ILIKE"
		if input.PaginationParams.CaseSensitiveSearch {
			matchType = "LIKE"
		}
		for i, field := range input.AllowedSearchFields {
			if i == 0 {
				query = query.Where(fmt.Sprintf("%s %s ?", field, matchType), searchClause)
			} else {
				query = query.Or(fmt.Sprintf("%s %s ?", field, matchType), searchClause)
			}
		}

	}

	filter := input.PaginationParams.Filter
	if len(filter) > 0 && len(input.AllowedFilterFields) > 0 {
		for _, f := range filter {
			if slices.Contains(input.AllowedFilterFields, f.Field) {
				query = query.Where(fmt.Sprintf("%s IN ?", f.Field), f.Value)
			}
		}
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, false, err
	}

	if input.PaginationParams.Sort.Field != "" && input.PaginationParams.Sort.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", input.PaginationParams.Sort.Field, input.PaginationParams.Sort.Order))
		sortApplied = true
	}

	if input.PaginationParams.Limit > 0 {
		query = query.Offset(input.PaginationParams.Offset).Limit(input.PaginationParams.Limit)
	}

	return query, totalRecords, sortApplied, nil
}
