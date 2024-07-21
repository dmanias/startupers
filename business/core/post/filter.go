package post

import (
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

type QueryFilter struct {
	ID               *uuid.UUID `validate:"omitempty"`
	IdeaID           *uuid.UUID `validate:"omitempty"`
	AuthorID         *uuid.UUID `validate:"omitempty"`
	OwnerType        *string    `validate:"omitempty"`
	StartCreatedDate *time.Time `validate:"omitempty"`
	EndCreatedDate   *time.Time `validate:"omitempty"`
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func (qf *QueryFilter) WithPostID(postID uuid.UUID) {
	qf.ID = &postID
}

func (qf *QueryFilter) WithIdeaID(ideaID uuid.UUID) {
	qf.IdeaID = &ideaID
}

func (qf *QueryFilter) WithAuthorID(authorID uuid.UUID) {
	qf.AuthorID = &authorID
}

func (qf *QueryFilter) WithOwnerType(ownerType string) {
	qf.OwnerType = &ownerType
}

func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
