package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Damoz1606/ministock-backend/db/sqlc"
	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

type categoryRepository struct {
	db *sqlc.Queries
}

func NewCategoryRepository(db *sqlc.Queries) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Save(ctx context.Context, category *domain.Category) error {
	_, err := r.db.CreateCategory(ctx, sqlc.CreateCategoryParams{
		ID:        pgtype.UUID{Bytes: [16]byte(category.ID), Valid: true},
		Name:      category.Name.String(),
		CreatedAt: category.CreatedAt.Unix(),
		UpdatedAt: category.UpdatedAt.Unix(),
	})
	return err
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	row, err := r.db.FindCategoryByID(ctx, pgtype.UUID{Bytes: [16]byte(id), Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return toDomainCategory(row), nil
}

func (r *categoryRepository) FindAll(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
	offset := (params.Page - 1) * params.PageSize

	rows, err := r.db.FindAllCategories(ctx, sqlc.FindAllCategoriesParams{
		Column1: params.NameFilter,
		Column2: params.OrderBy,
		Column3: params.OrderDir,
		Limit:   int32(params.PageSize),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	count, err := r.db.CountCategories(ctx, params.NameFilter)
	if err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(rows))
	for i, row := range rows {
		categories[i] = toDomainCategory(row)
	}

	return &domain.FindAllResult{
		Categories: categories,
		Total:      count,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	_, err := r.db.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		Name:      category.Name.String(),
		UpdatedAt: category.UpdatedAt.Unix(),
		ID:        pgtype.UUID{Bytes: [16]byte(category.ID), Valid: true},
	})
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteCategory(ctx, pgtype.UUID{Bytes: [16]byte(id), Valid: true})
}

func (r *categoryRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	exists, err := r.db.ExistsCategoryByName(ctx, sqlc.ExistsCategoryByNameParams{
		Name: name,
		ID:   pgtype.UUID{Bytes: [16]byte(excludeID), Valid: true},
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func toDomainCategory(row sqlc.Category) *domain.Category {
	return &domain.Category{
		ID:        uuid.UUID(row.ID.Bytes),
		Name:      domain.CategoryName(row.Name),
		CreatedAt: time.Unix(row.CreatedAt, 0),
		UpdatedAt: time.Unix(row.UpdatedAt, 0),
	}
}
