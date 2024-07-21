package challengegrp

import (
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/core/challenge"
	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

type AppChallenge struct {
	ID          string `json:"id"`
	IdeaID      string `json:"ideaID"`
	ModeratorID string `json:"moderatorID"`
	Answer      string `json:"answer"`
	PhotoURL    string `json:"photoURL"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppChallenge(challenge challenge.Challenge) AppChallenge {
	return AppChallenge{
		ID:          challenge.ID.String(),
		IdeaID:      challenge.IdeaID.String(),
		ModeratorID: challenge.ModeratorID.String(),
		Answer:      challenge.Answer,
		PhotoURL:    challenge.PhotoURL,
		DateCreated: challenge.DateCreated.Format(time.RFC3339),
		DateUpdated: challenge.DateUpdated.Format(time.RFC3339),
	}
}

type AppNewChallenge struct {
	IdeaID      string `json:"ideaID" validate:"required"`
	ModeratorID string `json:"moderatorID" validate:"required"`
	Answer      string `json:"answer"`
	PhotoURL    string `json:"photoURL"`
}

func toCoreNewChallenge(app AppNewChallenge) (challenge.NewChallenge, error) {
	ideaID, err := uuid.Parse(app.IdeaID)
	if err != nil {
		return challenge.NewChallenge{}, fmt.Errorf("parsing ideaID: %w", err)
	}

	moderatorID, err := uuid.Parse(app.ModeratorID)
	if err != nil {
		return challenge.NewChallenge{}, fmt.Errorf("parsing moderatorID: %w", err)
	}

	nc := challenge.NewChallenge{
		IdeaID:      ideaID,
		ModeratorID: moderatorID,
		Answer:      app.Answer,
		PhotoURL:    app.PhotoURL,
	}

	return nc, nil
}

func (app AppNewChallenge) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

type AppUpdateChallenge struct {
	ID       *string `json:"id"`
	Answer   *string `json:"answer"`
	PhotoURL *string `json:"photoURL"`
}

func toCoreUpdateChallenge(app AppUpdateChallenge) (challenge.UpdateChallenge, error) {
	var id uuid.UUID
	if app.ID != nil {
		var err error
		id, err = uuid.Parse(*app.ID)
		if err != nil {
			return challenge.UpdateChallenge{}, fmt.Errorf("parsing ID: %w", err)
		}
	}

	uc := challenge.UpdateChallenge{
		ID:       &id,
		Answer:   app.Answer,
		PhotoURL: app.PhotoURL,
	}

	return uc, nil
}

func (app AppUpdateChallenge) Validate() error {
	if err := validate.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
