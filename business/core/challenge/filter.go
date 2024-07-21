package challenge

import (
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

type QueryFilter struct {
	ID               *uuid.UUID `validate:"omitempty"`
	IdeaID           *uuid.UUID `validate:"omitempty"`
	ModeratorID      *uuid.UUID `validate:"omitempty"`
	StartCreatedDate *time.Time `validate:"omitempty"`
	EndCreatedDate   *time.Time `validate:"omitempty"`
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func (qf *QueryFilter) WithChallengeID(challengeID uuid.UUID) {
	qf.ID = &challengeID
}

func (qf *QueryFilter) WithIdeaID(ideaID uuid.UUID) {
	qf.IdeaID = &ideaID
}

func (qf *QueryFilter) WithModeratorID(moderatorID uuid.UUID) {
	qf.ModeratorID = &moderatorID
}

func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
