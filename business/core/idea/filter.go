package idea

import (
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

type QueryFilter struct {
	ID               *uuid.UUID `validate:"omitempty"`
	UserID           *uuid.UUID `validate:"omitempty"`
	Title            *string    `validate:"omitempty,min=3"`
	Category         *string    `validate:"omitempty"`
	Tag              *string    `validate:"omitempty"`
	StartCreatedDate *time.Time `validate:"omitempty"`
	EndCreatedDate   *time.Time `validate:"omitempty"`
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func (qf *QueryFilter) WithIdeaID(ideaID uuid.UUID) {
	qf.ID = &ideaID
}

func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.UserID = &userID
}

func (qf *QueryFilter) WithTitle(title string) {
	qf.Title = &title
}

func (qf *QueryFilter) WithCategory(category string) {
	qf.Category = &category
}

func (qf *QueryFilter) WithTag(tag string) {
	qf.Tag = &tag
}

func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
