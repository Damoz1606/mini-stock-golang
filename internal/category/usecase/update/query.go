package update

import "github.com/google/uuid"

// Command represents the input for updating a category.
type Command struct {
	CategoryID uuid.UUID
	Name       string
}
