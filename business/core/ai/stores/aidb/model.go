package aidb

import (
	"github.com/dmanias/startupers/business/core/ai"
	"github.com/google/uuid"
	"time"
)

// dbModerator represent the structure we need for moving data
// between the app and the database.
type dbAi struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Query       string    `db:"query"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBAi(mdr ai.Ai) dbAi {

	return dbAi{
		ID:          mdr.ID,
		Name:        mdr.Name,
		Query:       mdr.Query,
		DateCreated: mdr.DateCreated.UTC(),
		DateUpdated: mdr.DateUpdated.UTC(),
	}
}

func toCoreAi(dbMdr dbAi) ai.Ai {

	mdr := ai.Ai{
		ID:          dbMdr.ID,
		Name:        dbMdr.Name,
		Query:       dbMdr.Query,
		DateCreated: dbMdr.DateCreated.In(time.Local),
		DateUpdated: dbMdr.DateUpdated.In(time.Local),
	}

	return mdr
}

func toCoreAiSlice(dbAis []dbAi) []ai.Ai {
	mdrs := make([]ai.Ai, len(dbAis))
	for i, dbMdr := range dbAis {
		mdrs[i] = toCoreAi(dbMdr)
	}
	return mdrs
}
