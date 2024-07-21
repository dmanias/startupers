package aigrp

import (
	"fmt"
	"github.com/dmanias/startupers/business/core/ai"
	"time"

	"github.com/dmanias/startupers/business/sys/validate"
)

// AppUser represents information about an individual ai.
type AppAi struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Query       string `json:"query"`
	UserID      string `json:"userID"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppAi(ai ai.Ai) AppAi {
	return AppAi{
		ID:          ai.ID.String(),
		Name:        ai.Name,
		Query:       ai.Query,
		UserID:      ai.UserID.String(),
		DateCreated: ai.DateCreated.Format(time.RFC3339),
		DateUpdated: ai.DateUpdated.Format(time.RFC3339),
	}
}

// =============================================================================

// AppNewUser contains information needed to create a new ai.
type AppNewAi struct {
	Name  string `json:"name" validate:"required"`
	Query string `json:"query" validate:"required"`
}

func toCoreNewAi(app AppNewAi) (ai.NewAi, error) {

	ai := ai.NewAi{
		Name:  app.Name,
		Query: app.Query,
	}

	return ai, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewAi) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

// =============================================================================

// AppUpdateUser contains information needed to update a user.
type AppUpdateAi struct {
	Name  *string `json:"name"`
	Query *string `json:"query"`
}

func toCoreUpdateAi(app AppUpdateAi) (ai.UpdateAi, error) {
	nu := ai.UpdateAi{
		Name:  app.Name,
		Query: app.Query,
	}

	return nu, nil
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateAi) Validate() error {
	if err := validate.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// =============================================================================

// AppSummary represents information about an individual ai
type AppSummary struct {
	AiID       string `json:"aiID"`
	AiName     string `json:"aiName"`
	TotalCount int    `json:"totalCount"`
}
