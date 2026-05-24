package delete

import (
	"context"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

type DeleteHandler struct {
	repo domain.CategoryRepository
}

func NewDeleteHandler(repo domain.CategoryRepository) *DeleteHandler {
	return &DeleteHandler{repo: repo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd Command) error {
	_, err := h.repo.FindByID(ctx, cmd.CategoryID)
	if err != nil {
		return err
	}

	return h.repo.Delete(ctx, cmd.CategoryID)
}
