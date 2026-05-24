package readAll

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Damoz1606/ministock-backend/internal/category/domain"
)

type mockCategoryRepository struct {
	saveFunc         func(ctx context.Context, category *domain.Category) error
	findByIDFunc     func(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	findAllFunc      func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error)
	updateFunc       func(ctx context.Context, category *domain.Category) error
	deleteFunc       func(ctx context.Context, id uuid.UUID) error
	existsByNameFunc func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error)
}

func (m *mockCategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	return m.saveFunc(ctx, category)
}

func (m *mockCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockCategoryRepository) FindAll(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
	return m.findAllFunc(ctx, params)
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return m.updateFunc(ctx, category)
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFunc(ctx, id)
}

func (m *mockCategoryRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	return m.existsByNameFunc(ctx, name, excludeID)
}

func TestReadAllHandler_Handle(t *testing.T) {
	t.Run("Should return categories and total count when categories exist", func(t *testing.T) {
		now := time.Now()
		categories := []*domain.Category{
			{ID: uuid.New(), Name: "Electronics", CreatedAt: now, UpdatedAt: now},
			{ID: uuid.New(), Name: "Books", CreatedAt: now, UpdatedAt: now},
		}
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				return &domain.FindAllResult{
					Categories: categories,
					Total:      2,
					Page:       1,
					PageSize:   20,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 20}
		result, total, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
		assert.Equal(t, "Electronics", string(result[0].Name))
		assert.Equal(t, "Books", string(result[1].Name))
	})

	t.Run("Should return empty slice when no categories exist", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				return &domain.FindAllResult{
					Categories: []*domain.Category{},
					Total:      0,
					Page:       1,
					PageSize:   20,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 20}
		result, total, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})

	t.Run("Should default page to 1 when page is less than 1", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				assert.Equal(t, 1, params.Page)
				return &domain.FindAllResult{
					Categories: []*domain.Category{},
					Total:      0,
					Page:       1,
					PageSize:   20,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 0, PageSize: 20}
		_, _, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
	})

	t.Run("Should default pageSize to 20 when pageSize is less than 1", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				assert.Equal(t, 20, params.PageSize)
				return &domain.FindAllResult{
					Categories: []*domain.Category{},
					Total:      0,
					Page:       1,
					PageSize:   20,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 0}
		_, _, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
	})

	t.Run("Should cap pageSize at 100 when pageSize exceeds 100", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				assert.Equal(t, 100, params.PageSize)
				return &domain.FindAllResult{
					Categories: []*domain.Category{},
					Total:      0,
					Page:       1,
					PageSize:   100,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 200}
		_, _, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
	})

	t.Run("Should default orderBy to created_at when invalid", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				assert.Equal(t, "created_at", params.OrderBy)
				return &domain.FindAllResult{
					Categories: []*domain.Category{},
					Total:      0,
					Page:       1,
					PageSize:   20,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 20, OrderBy: "invalid"}
		_, _, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
	})

	t.Run("Should default orderDir to asc when invalid", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				assert.Equal(t, "asc", params.OrderDir)
				return &domain.FindAllResult{
					Categories: []*domain.Category{},
					Total:      0,
					Page:       1,
					PageSize:   20,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 20, OrderDir: "invalid"}
		_, _, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
	})

	t.Run("Should pass name filter to repository", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				assert.Equal(t, "Electronics", params.NameFilter)
				return &domain.FindAllResult{
					Categories: []*domain.Category{},
					Total:      0,
					Page:       1,
					PageSize:   20,
				}, nil
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 20, NameFilter: "Electronics"}
		_, _, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
	})

	t.Run("Should return error when FindAll fails", func(t *testing.T) {
		expectedErr := errors.New("database query error")
		mockRepo := &mockCategoryRepository{
			findAllFunc: func(ctx context.Context, params domain.FindAllParams) (*domain.FindAllResult, error) {
				return nil, expectedErr
			},
		}
		handler := NewReadAllHandler(mockRepo)

		q := Query{Page: 1, PageSize: 20}
		result, total, err := handler.Handle(context.Background(), q)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
	})
}
