package moderator

import (
	"github.com/google/uuid"
	"time"
)

// Moderator represents an individual product.
type Moderator struct {
	ID          uuid.UUID
	Name        string
	Instruction string
	DateCreated time.Time
	DateUpdated time.Time
}

// NewModerator is what we require from clients when adding a Product.
type NewModerator struct {
	Name        string
	Instruction string
}

// UpdateModerator defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateModerator struct {
	Name        *string
	Instruction *string
}
