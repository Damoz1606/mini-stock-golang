package delete

import (
	"context"
	"errors"
	"testing"

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

func TestDeleteHandler_Handle(t *testing.T) {
	t.Run("Should return nil when category is deleted successfully", func(t *testing.T) {
		id := uuid.New()
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return &domain.Category{ID: foundID}, nil
			},
			deleteFunc: func(ctx context.Context, deleteID uuid.UUID) error {
				return nil
			},
		}
		handler := NewDeleteHandler(mockRepo)

		cmd := Command{CategoryID: id}
		err := handler.Handle(context.Background(), cmd)

		assert.NoError(t, err)
	})

	t.Run("Should return ErrNotFound when category does not exist", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
				return nil, domain.ErrNotFound
			},
		}
		handler := NewDeleteHandler(mockRepo)

		cmd := Command{CategoryID: uuid.New()}
		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrNotFound, err)
	})

	t.Run("Should return error when FindByID fails with database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
				return nil, expectedErr
			},
		}
		handler := NewDeleteHandler(mockRepo)

		cmd := Command{CategoryID: uuid.New()}
		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when Delete fails", func(t *testing.T) {
		expectedErr := errors.New("database delete error")
		id := uuid.New()
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return &domain.Category{ID: foundID}, nil
			},
			deleteFunc: func(ctx context.Context, deleteID uuid.UUID) error {
				return expectedErr
			},
		}
		handler := NewDeleteHandler(mockRepo)

		cmd := Command{CategoryID: id}
		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
