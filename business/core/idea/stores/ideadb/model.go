package ideadb

import (
	"time"

	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/sys/database/pgx/dbarray"
	"github.com/google/uuid"
)

type dbIdea struct {
	ID            uuid.UUID      `db:"id"`
	UserID        uuid.UUID      `db:"user_id"`
	Title         string         `db:"title"`
	Description   string         `db:"description"`
	Category      string         `db:"category"`
	Tags          dbarray.String `db:"tags"`
	Privacy       string         `db:"privacy"`
	Collaborators dbarray.UUID   `db:"collaborators"`
	AvatarURL     string         `db:"avatar_url"`
	Stage         string         `db:"stage"`
	Inspiration   string         `db:"inspiration"`
	DateCreated   time.Time      `db:"date_created"`
	DateUpdated   time.Time      `db:"date_updated"`
}

func toDBIdea(idea idea.Idea) dbIdea {
	return dbIdea{
		ID:            idea.ID,
		UserID:        idea.UserID,
		Title:         idea.Title,
		Description:   idea.Description,
		Category:      idea.Category,
		Tags:          idea.Tags,
		Privacy:       idea.Privacy,
		Collaborators: idea.Collaborators,
		AvatarURL:     idea.AvatarURL,
		Stage:         idea.Stage,
		Inspiration:   idea.Inspiration,
		DateCreated:   idea.DateCreated.UTC(),
		DateUpdated:   idea.DateUpdated.UTC(),
	}
}

func toCoreIdea(dbIdea dbIdea) idea.Idea {
	return idea.Idea{
		ID:            dbIdea.ID,
		UserID:        dbIdea.UserID,
		Title:         dbIdea.Title,
		Description:   dbIdea.Description,
		Category:      dbIdea.Category,
		Tags:          dbIdea.Tags,
		Privacy:       dbIdea.Privacy,
		Collaborators: dbIdea.Collaborators,
		AvatarURL:     dbIdea.AvatarURL,
		Stage:         dbIdea.Stage,
		Inspiration:   dbIdea.Inspiration,
		DateCreated:   dbIdea.DateCreated.In(time.Local),
		DateUpdated:   dbIdea.DateUpdated.In(time.Local),
	}
}

func toCoreIdeaSlice(dbIdeas []dbIdea) []idea.Idea {
	ideas := make([]idea.Idea, len(dbIdeas))
	for i, dbIdea := range dbIdeas {
		ideas[i] = toCoreIdea(dbIdea)
	}
	return ideas
}
