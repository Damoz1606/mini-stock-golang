package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

type CreateHandler struct {
	repo domain.CategoryRepository
}

func NewCreateHandler(repo domain.CategoryRepository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd Command) (*domain.Category, error) {
	name, err := domain.NewCategoryName(cmd.Name)
	if err != nil {
		return nil, domain.ErrInvalidName
	}

	exists, err := h.repo.ExistsByName(ctx, name.String(), uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrDuplicateName
	}

	now := time.Now()
	category, err := domain.NewCategory(uuid.New(), name, now, now)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
