package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CategoryName string

func NewCategoryName(v string) (CategoryName, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", ErrInvalidName
	}
	if len(v) > 255 {
		return "", ErrInvalidName
	}
	return CategoryName(v), nil
}

func (cn CategoryName) String() string {
	return string(cn)
}

type Category struct {
	ID        uuid.UUID
	Name      CategoryName
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCategory(id uuid.UUID, name CategoryName, createdAt, updatedAt time.Time) (*Category, error) {
	if id == uuid.Nil {
		return nil, errors.New("category id is required")
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	return &Category{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
