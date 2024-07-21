package idea

import (
	"time"

	"github.com/google/uuid"
)

type Idea struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Title         string
	Description   string
	Category      string
	Tags          []string
	Privacy       string
	Collaborators []uuid.UUID
	AvatarURL     string
	Stage         string
	Inspiration   string
	DateCreated   time.Time
	DateUpdated   time.Time
}

type NewIdea struct {
	UserID        uuid.UUID
	Title         string
	Description   string
	Category      string
	Tags          []string
	Privacy       string
	Collaborators []uuid.UUID
	AvatarURL     string
	Stage         string
	Inspiration   string
}

type UpdateIdea struct {
	ID            *uuid.UUID
	Title         *string
	Description   *string
	Category      *string
	Tags          []string
	Privacy       *string
	Collaborators []uuid.UUID
	AvatarURL     *string
	Stage         *string
	Inspiration   *string
}
