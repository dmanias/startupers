package ai

import (
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	ID               *uuid.UUID `validate:"omitempty"`
	Name             *string    `validate:"omitempty,min=3"`
	Query            *string    `validate:"omitempty,min=20"`
	StartCreatedDate *time.Time `validate:"omitempty"`
	EndCreatedDate   *time.Time `validate:"omitempty"`
}

// Validate checks the data in the model is considered clean.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithModeratorID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithAiID(aiID uuid.UUID) {
	qf.ID = &aiID
}

// WithName sets the Name field of the QueryFilter value.
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithCost sets the Cost field of the QueryFilter value.
func (qf *QueryFilter) WithInstruction(query string) {
	qf.Query = &query
}

// WithStartDateCreated sets the DateCreated field of the QueryFilter value.
func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

// WithEndCreatedDate sets the DateCreated field of the QueryFilter value.
func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
