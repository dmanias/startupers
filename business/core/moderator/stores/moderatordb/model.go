package moderatordb

import (
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/google/uuid"
	"time"
)

// dbModerator represent the structure we need for moving data
// between the app and the database.
type dbModerator struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Instruction string    `db:"instruction"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBModerator(mdr moderator.Moderator) dbModerator {

	return dbModerator{
		ID:          mdr.ID,
		Name:        mdr.Name,
		Instruction: mdr.Instruction,
		DateCreated: mdr.DateCreated.UTC(),
		DateUpdated: mdr.DateUpdated.UTC(),
	}
}

func toCoreModerator(dbMdr dbModerator) moderator.Moderator {

	mdr := moderator.Moderator{
		ID:          dbMdr.ID,
		Name:        dbMdr.Name,
		Instruction: dbMdr.Instruction,
		DateCreated: dbMdr.DateCreated.In(time.Local),
		DateUpdated: dbMdr.DateUpdated.In(time.Local),
	}

	return mdr
}

func toCoreModeratorSlice(dbModerators []dbModerator) []moderator.Moderator {
	mdrs := make([]moderator.Moderator, len(dbModerators))
	for i, dbMdr := range dbModerators {
		mdrs[i] = toCoreModerator(dbMdr)
	}
	return mdrs
}
