package response

import (
	"github.com/google/uuid"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

// categoryResponse is the shared response type for create, readOne, and update operations.
type categoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

// ToCategoryResponse converts a domain Category to a response struct.
func ToCategoryResponse(category *domain.Category) categoryResponse {
	return categoryResponse{
		ID:        category.ID,
		Name:      category.Name.String(),
		CreatedAt: category.CreatedAt.Unix(),
		UpdatedAt: category.UpdatedAt.Unix(),
	}
}
