package moderationgrp

import (
	"fmt"
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/dmanias/startupers/business/sys/validate"
	"time"
)

// AppUser represents information about an individual ai.
type AppModerator struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Instruction string `json:"instruction"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppModeration(mdr moderator.Moderator) AppModerator {
	return AppModerator{
		ID:          mdr.ID.String(),
		Name:        mdr.Name,
		Instruction: mdr.Instruction,
		DateCreated: mdr.DateCreated.Format(time.RFC3339),
		DateUpdated: mdr.DateUpdated.Format(time.RFC3339),
	}
}

// =============================================================================

// AppNewUser contains information needed to create a new ai.
type AppNewModerator struct {
	Name        string `json:"name" validate:"required"`
	Instruction string `json:"Instruction" validate:"required"`
}

func toCoreNewModerator(app AppNewModerator) (moderator.NewModerator, error) {

	mdr := moderator.NewModerator{
		Name:        app.Name,
		Instruction: app.Instruction,
	}

	return mdr, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewModerator) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

// =============================================================================

// AppUpdateUser contains information needed to update a user.
type AppUpdateModerator struct {
	Name        *string `json:"name"`
	Instruction *string `json:"Instruction"`
}

func toCoreUpdateModerator(app AppUpdateModerator) (moderator.UpdateModerator, error) {
	nu := moderator.UpdateModerator{
		Name:        app.Name,
		Instruction: app.Instruction,
	}

	return nu, nil
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateModerator) Validate() error {
	if err := validate.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
