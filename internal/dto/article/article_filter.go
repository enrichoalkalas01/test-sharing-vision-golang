package dto

import "github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/common/filter"

type ArticleFilter struct {
	filter.BaseFilter[ArticleFilterFields]
}

type ArticleFilterFields struct {
	Category string `query:"category"`
	Status   string `query:"status"`
}

func NewArticleFilter() *ArticleFilter {
	return &ArticleFilter{
		BaseFilter: filter.BaseFilter[ArticleFilterFields]{
			Page:      1,
			Limit:     10,
			SortBy:    "created_date",
			SortOrder: "desc",
		},
	}
}

func (af *ArticleFilter) Validate() error {
	// Validasi pagination
	if err := af.ValidatePagination(); err != nil {
		return err
	}

	// Validasi status jika ada
	if af.Filters.Status != "" {
		validStatuses := map[string]bool{
			"Publish": true,
			"Draft":   true,
			"Thrash":  true,
		}
		if !validStatuses[af.Filters.Status] {
			return ErrInvalidFilterStatus
		}
	}

	return nil
}

func (af *ArticleFilter) GetCategory() string {
	return af.Filters.Category
}

func (af *ArticleFilter) GetStatus() string {
	return af.Filters.Status
}

func (af *ArticleFilter) HasCategory() bool {
	return af.GetCategory() != ""
}

func (af *ArticleFilter) HasStatus() bool {
	return af.GetStatus() != ""
}

func (af *ArticleFilter) BuildQueryConditions() []filter.QueryCondition {
	var conditions []filter.QueryCondition

	if af.HasCategory() {
		conditions = append(conditions, filter.QueryCondition{
			Field:    "category",
			Operator: "=",
			Value:    af.GetCategory(),
		})
	}

	if af.HasStatus() {
		conditions = append(conditions, filter.QueryCondition{
			Field:    "status",
			Operator: "=",
			Value:    af.GetStatus(),
		})
	}

	if af.GetSearch() != "" {
		conditions = append(conditions, filter.QueryCondition{
			Field:    "title",
			Operator: "LIKE",
			Value:    "%" + af.GetSearch() + "%",
		})
	}

	return conditions
}
