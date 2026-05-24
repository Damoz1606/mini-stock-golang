package create

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

func TestCreateHandler_Handle(t *testing.T) {
	t.Run("Should return created category when name is valid and unique", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return false, nil
			},
			saveFunc: func(ctx context.Context, category *domain.Category) error {
				return nil
			},
		}
		handler := NewCreateHandler(mockRepo)

		cmd := Command{Name: "Electronics"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, domain.CategoryName("Electronics"), category.Name)
		assert.NotEqual(t, uuid.Nil, category.ID)
	})

	t.Run("Should return ErrInvalidName when name is empty", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{}
		handler := NewCreateHandler(mockRepo)

		cmd := Command{Name: ""}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidName, err)
		assert.Nil(t, category)
	})

	t.Run("Should return ErrInvalidName when name is whitespace only", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{}
		handler := NewCreateHandler(mockRepo)

		cmd := Command{Name: "   "}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidName, err)
		assert.Nil(t, category)
	})

	t.Run("Should return ErrDuplicateName when name already exists", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},
		}
		handler := NewCreateHandler(mockRepo)

		cmd := Command{Name: "Electronics"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrDuplicateName, err)
		assert.Nil(t, category)
	})

	t.Run("Should return error when ExistsByName fails", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		mockRepo := &mockCategoryRepository{
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return false, expectedErr
			},
		}
		handler := NewCreateHandler(mockRepo)

		cmd := Command{Name: "Electronics"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, category)
	})

	t.Run("Should return error when Save fails", func(t *testing.T) {
		expectedErr := errors.New("database write error")
		mockRepo := &mockCategoryRepository{
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return false, nil
			},
			saveFunc: func(ctx context.Context, category *domain.Category) error {
				return expectedErr
			},
		}
		handler := NewCreateHandler(mockRepo)

		cmd := Command{Name: "Electronics"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, category)
	})
}
