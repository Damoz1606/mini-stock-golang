package domain

import (
	"context"

	"github.com/google/uuid"
)

type FindAllParams struct {
	Page       int
	PageSize   int
	OrderBy    string
	OrderDir   string
	NameFilter string
}

type FindAllResult struct {
	Categories []*Category
	Total      int64
	Page       int
	PageSize   int
}

type CategoryRepository interface {
	Save(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	FindAll(ctx context.Context, params FindAllParams) (*FindAllResult, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error)
}
