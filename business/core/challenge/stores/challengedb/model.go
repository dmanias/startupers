package challengedb

import (
	"github.com/dmanias/startupers/business/core/challenge"
	"time"

	"github.com/google/uuid"
)

type dbChallenge struct {
	ID          uuid.UUID `db:"id"`
	IdeaID      uuid.UUID `db:"idea_id"`
	ModeratorID uuid.UUID `db:"moderator_id"`
	Answer      string    `db:"answer"`
	PhotoURL    string    `db:"photo_url"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBChallenge(challenge challenge.Challenge) dbChallenge {
	return dbChallenge{
		ID:          challenge.ID,
		IdeaID:      challenge.IdeaID,
		ModeratorID: challenge.ModeratorID,
		Answer:      challenge.Answer,
		PhotoURL:    challenge.PhotoURL,
		DateCreated: challenge.DateCreated.UTC(),
		DateUpdated: challenge.DateUpdated.UTC(),
	}
}

func toCoreChallenge(dbChallenge dbChallenge) challenge.Challenge {
	return challenge.Challenge{
		ID:          dbChallenge.ID,
		IdeaID:      dbChallenge.IdeaID,
		ModeratorID: dbChallenge.ModeratorID,
		Answer:      dbChallenge.Answer,
		PhotoURL:    dbChallenge.PhotoURL,
		DateCreated: dbChallenge.DateCreated.In(time.Local),
		DateUpdated: dbChallenge.DateUpdated.In(time.Local),
	}
}

func toCoreChallengeSlice(dbChallenges []dbChallenge) []challenge.Challenge {
	challenges := make([]challenge.Challenge, len(dbChallenges))
	for i, dbChallenge := range dbChallenges {
		challenges[i] = toCoreChallenge(dbChallenge)
	}
	return challenges
}
