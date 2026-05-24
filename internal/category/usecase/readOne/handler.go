package readOne

import (
	"context"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

type ReadOneHandler struct {
	repo domain.CategoryRepository
}

func NewReadOneHandler(repo domain.CategoryRepository) *ReadOneHandler {
	return &ReadOneHandler{repo: repo}
}

func (h *ReadOneHandler) Handle(ctx context.Context, q Query) (*domain.Category, error) {
	return h.repo.FindByID(ctx, q.CategoryID)
}
