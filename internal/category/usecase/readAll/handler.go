package readAll

import (
	"context"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

type ReadAllHandler struct {
	repo domain.CategoryRepository
}

func NewReadAllHandler(repo domain.CategoryRepository) *ReadAllHandler {
	return &ReadAllHandler{repo: repo}
}

func (h *ReadAllHandler) Handle(ctx context.Context, q Query) ([]*domain.Category, int64, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}

	validOrderBy := q.OrderBy
	if validOrderBy != "name" && validOrderBy != "created_at" {
		validOrderBy = "created_at"
	}

	validOrderDir := q.OrderDir
	if validOrderDir != "asc" && validOrderDir != "desc" {
		validOrderDir = "asc"
	}

	params := domain.FindAllParams{
		Page:       q.Page,
		PageSize:   q.PageSize,
		OrderBy:    validOrderBy,
		OrderDir:   validOrderDir,
		NameFilter: q.NameFilter,
	}

	result, err := h.repo.FindAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return result.Categories, result.Total, nil
}
