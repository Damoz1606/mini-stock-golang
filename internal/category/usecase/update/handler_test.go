package update

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

func TestUpdateHandler_Handle(t *testing.T) {
	t.Run("Should return updated category when name is valid and unique", func(t *testing.T) {
		id := uuid.New()
		now := time.Now()
		existingCategory := &domain.Category{
			ID:        id,
			Name:      "Old Name",
			CreatedAt: now,
			UpdatedAt: now,
		}
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return existingCategory, nil
			},
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return false, nil
			},
			updateFunc: func(ctx context.Context, category *domain.Category) error {
				return nil
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: id, Name: "New Name"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, domain.CategoryName("New Name"), category.Name)
		assert.Equal(t, id, category.ID)
	})

	t.Run("Should return ErrNotFound when category does not exist", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
				return nil, domain.ErrNotFound
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: uuid.New(), Name: "New Name"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrNotFound, err)
		assert.Nil(t, category)
	})

	t.Run("Should return ErrInvalidName when name is empty", func(t *testing.T) {
		id := uuid.New()
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return &domain.Category{ID: foundID}, nil
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: id, Name: ""}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidName, err)
		assert.Nil(t, category)
	})

	t.Run("Should return ErrInvalidName when name is whitespace only", func(t *testing.T) {
		id := uuid.New()
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return &domain.Category{ID: foundID}, nil
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: id, Name: "   "}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidName, err)
		assert.Nil(t, category)
	})

	t.Run("Should return ErrDuplicateName when name already exists on another category", func(t *testing.T) {
		id := uuid.New()
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return &domain.Category{ID: foundID}, nil
			},
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return true, nil
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: id, Name: "Existing Name"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrDuplicateName, err)
		assert.Nil(t, category)
	})

	t.Run("Should return error when FindByID fails with database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
				return nil, expectedErr
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: uuid.New(), Name: "New Name"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, category)
	})

	t.Run("Should return error when ExistsByName fails", func(t *testing.T) {
		expectedErr := errors.New("database query error")
		id := uuid.New()
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return &domain.Category{ID: foundID}, nil
			},
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return false, expectedErr
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: id, Name: "New Name"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, category)
	})

	t.Run("Should return error when Update fails", func(t *testing.T) {
		expectedErr := errors.New("database write error")
		id := uuid.New()
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return &domain.Category{ID: foundID}, nil
			},
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return false, nil
			},
			updateFunc: func(ctx context.Context, category *domain.Category) error {
				return expectedErr
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: id, Name: "New Name"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, category)
	})

	t.Run("Should update UpdatedAt timestamp on successful update", func(t *testing.T) {
		id := uuid.New()
		now := time.Now()
		existingCategory := &domain.Category{
			ID:        id,
			Name:      "Old Name",
			CreatedAt: now,
			UpdatedAt: now,
		}
		mockRepo := &mockCategoryRepository{
			findByIDFunc: func(ctx context.Context, foundID uuid.UUID) (*domain.Category, error) {
				return existingCategory, nil
			},
			existsByNameFunc: func(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
				return false, nil
			},
			updateFunc: func(ctx context.Context, category *domain.Category) error {
				assert.True(t, category.UpdatedAt.After(now), "UpdatedAt should be updated to current time")
				return nil
			},
		}
		handler := NewUpdateHandler(mockRepo)

		cmd := Command{CategoryID: id, Name: "New Name"}
		category, err := handler.Handle(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.True(t, category.UpdatedAt.After(now), "UpdatedAt should be after the original timestamp")
	})
}
