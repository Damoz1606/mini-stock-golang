package response

import (
	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

// listCategoriesResponse is the response type for the list categories operation.
type listCategoriesResponse struct {
	Categories []categoryResponse `json:"categories"`
	TotalCount int64              `json:"totalCount"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	TotalPages int                `json:"totalPages"`
}

// ToListCategoriesResponse converts a slice of domain Categories to a paginated response.
func ToListCategoriesResponse(categories []*domain.Category, total int64, page, pageSize int) listCategoriesResponse {
	items := make([]categoryResponse, len(categories))
	for i, c := range categories {
		items[i] = ToCategoryResponse(c)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return listCategoriesResponse{
		Categories: items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
