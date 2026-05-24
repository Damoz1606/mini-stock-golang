package domain

import "errors"

var (
	// ErrNotFound is returned when a category is not found by the given ID.
	ErrNotFound = errors.New("category not found")

	// ErrDuplicateName is returned when a category name already exists.
	ErrDuplicateName = errors.New("category name already exists")

	// ErrInvalidName is returned when a category name fails validation.
	ErrInvalidName = errors.New("invalid category name")
)
