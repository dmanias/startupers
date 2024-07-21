package challenge

import (
	"time"

	"github.com/google/uuid"
)

type Challenge struct {
	ID          uuid.UUID
	IdeaID      uuid.UUID
	ModeratorID uuid.UUID
	Answer      string
	PhotoURL    string
	DateCreated time.Time
	DateUpdated time.Time
}

type NewChallenge struct {
	IdeaID      uuid.UUID
	ModeratorID uuid.UUID
	Answer      string
	PhotoURL    string
}

type UpdateChallenge struct {
	ID       *uuid.UUID
	Answer   *string
	PhotoURL *string
}
