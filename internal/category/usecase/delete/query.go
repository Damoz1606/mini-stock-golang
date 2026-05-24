package delete

import "github.com/google/uuid"

type Command struct {
	CategoryID uuid.UUID
}
