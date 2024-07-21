package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID
	IdeaID      uuid.UUID
	AuthorID    uuid.UUID
	Content     string
	OwnerType   string
	DateCreated time.Time
	DateUpdated time.Time
}

type NewPost struct {
	IdeaID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	OwnerType string
}

type UpdatePost struct {
	ID      *uuid.UUID
	Content *string
}
