package readOne

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

func TestReadOneHandler_Handle(t *testing.T) {
	t.Run("Should return category when category exists", func(t *testing.T) {
		id := uuid.New()
		now := time.Now()
		expectedCategory := &domain.Category{
			ID:        id,
			Name:      "Electronics",
			CreatedAt: now,
			UpdatedAt: now,
		}
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return expectedCategory, nil
			},
		}
		handler := NewReadOneHandler(mockRepo)

		q := Query{CategoryID: id}
		category, err := handler.Handle(context.Background(), q)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, id, category.ID)
		assert.Equal(t, domain.CategoryName("Electronics"), category.Name)
	})

	t.Run("Should return ErrNotFound when category does not exist", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
				return nil, domain.ErrNotFound
			},
		}
		handler := NewReadOneHandler(mockRepo)

		q := Query{CategoryID: uuid.New()}
		category, err := handler.Handle(context.Background(), q)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrNotFound, err)
		assert.Nil(t, category)
	})

	t.Run("Should return error when FindByID fails with database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
				return nil, expectedErr
			},
		}
		handler := NewReadOneHandler(mockRepo)

		q := Query{CategoryID: uuid.New()}
		category, err := handler.Handle(context.Background(), q)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, category)
	})
}
