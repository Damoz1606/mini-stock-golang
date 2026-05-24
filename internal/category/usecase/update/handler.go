package update

import (
	"context"
	"time"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

// UpdateHandler handles the update category usecase.
type UpdateHandler struct {
	repo domain.CategoryRepository
}

// NewUpdateHandler creates a new UpdateHandler.
func NewUpdateHandler(repo domain.CategoryRepository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

// Handle executes the update category usecase.
// It validates the new name, checks existence and uniqueness, then updates and saves.
func (h *UpdateHandler) Handle(ctx context.Context, cmd Command) (*domain.Category, error) {
	category, err := h.repo.FindByID(ctx, cmd.CategoryID)
	if err != nil {
		return nil, err
	}

	name, err := domain.NewCategoryName(cmd.Name)
	if err != nil {
		return nil, domain.ErrInvalidName
	}

	exists, err := h.repo.ExistsByName(ctx, name.String(), cmd.CategoryID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrDuplicateName
	}

	category.Name = name
	category.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
