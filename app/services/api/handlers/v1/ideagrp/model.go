package ideagrp

import (
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

// AppIdea represents information about an individual idea.
type AppIdea struct {
	ID            string   `json:"id"`
	UserID        string   `json:"userID"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	Privacy       string   `json:"privacy"`
	Collaborators []string `json:"collaborators"`
	AvatarURL     string   `json:"avatarURL"`
	Stage         string   `json:"stage"`
	Inspiration   string   `json:"inspiration"`
	DateCreated   string   `json:"dateCreated"`
	DateUpdated   string   `json:"dateUpdated"`
}

func toAppIdea(idea idea.Idea) AppIdea {
	collaborators := make([]string, len(idea.Collaborators))
	for i, collaborator := range idea.Collaborators {
		collaborators[i] = collaborator.String()
	}

	return AppIdea{
		ID:            idea.ID.String(),
		UserID:        idea.UserID.String(),
		Title:         idea.Title,
		Description:   idea.Description,
		Category:      idea.Category,
		Tags:          idea.Tags,
		Privacy:       idea.Privacy,
		Collaborators: collaborators,
		AvatarURL:     idea.AvatarURL,
		Stage:         idea.Stage,
		Inspiration:   idea.Inspiration,
		DateCreated:   idea.DateCreated.Format(time.RFC3339),
		DateUpdated:   idea.DateUpdated.Format(time.RFC3339),
	}
}

// AppNewIdea contains information needed to create a new idea.
type AppNewIdea struct {
	UserID        string   `json:"userID" validate:"required"`
	Title         string   `json:"title" validate:"required"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	Privacy       string   `json:"privacy"`
	Collaborators []string `json:"collaborators"`
	AvatarURL     string   `json:"avatarURL"`
	Stage         string   `json:"stage"`
	Inspiration   string   `json:"inspiration"`
}

func toCoreNewIdea(app AppNewIdea) (idea.NewIdea, error) {
	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		return idea.NewIdea{}, fmt.Errorf("parsing userID: %w", err)
	}

	collaborators := make([]uuid.UUID, len(app.Collaborators))
	for i, collaboratorStr := range app.Collaborators {
		collaborator, err := uuid.Parse(collaboratorStr)
		if err != nil {
			return idea.NewIdea{}, fmt.Errorf("parsing collaborator: %w", err)
		}
		collaborators[i] = collaborator
	}

	ni := idea.NewIdea{
		UserID:        userID,
		Title:         app.Title,
		Description:   app.Description,
		Category:      app.Category,
		Tags:          app.Tags,
		Privacy:       app.Privacy,
		Collaborators: collaborators,
		AvatarURL:     app.AvatarURL,
		Stage:         app.Stage,
		Inspiration:   app.Inspiration,
	}

	return ni, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewIdea) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

// AppUpdateIdea contains information needed to update an idea.
type AppUpdateIdea struct {
	ID            *string  `json:"id"`
	Title         *string  `json:"title"`
	Description   *string  `json:"description"`
	Category      *string  `json:"category"`
	Tags          []string `json:"tags"`
	Privacy       *string  `json:"privacy"`
	Collaborators []string `json:"collaborators"`
	AvatarURL     *string  `json:"avatarURL"`
	Stage         *string  `json:"stage"`
	Inspiration   *string  `json:"inspiration"`
}

func toCoreUpdateIdea(app AppUpdateIdea) (idea.UpdateIdea, error) {
	var id uuid.UUID
	if app.ID != nil {
		var err error
		id, err = uuid.Parse(*app.ID)
		if err != nil {
			return idea.UpdateIdea{}, fmt.Errorf("parsing ID: %w", err)
		}
	}
	var collaborators []uuid.UUID
	if app.Collaborators != nil {
		collaborators = make([]uuid.UUID, len(app.Collaborators))
		for i, collaboratorStr := range app.Collaborators {
			collaborator, err := uuid.Parse(collaboratorStr)
			if err != nil {
				return idea.UpdateIdea{}, fmt.Errorf("parsing collaborator: %w", err)
			}
			collaborators[i] = collaborator
		}
	}

	ui := idea.UpdateIdea{
		ID:            &id,
		Title:         app.Title,
		Description:   app.Description,
		Category:      app.Category,
		Tags:          app.Tags,
		Privacy:       app.Privacy,
		Collaborators: collaborators,
		AvatarURL:     app.AvatarURL,
		Stage:         app.Stage,
		Inspiration:   app.Inspiration,
	}

	return ui, nil
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateIdea) Validate() error {
	if err := validate.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
