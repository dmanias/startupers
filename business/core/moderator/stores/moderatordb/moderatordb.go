// Package challengedb contains user related CRUD functionality.
package moderatordb

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/dmanias/startupers/business/core/user"
	"github.com/dmanias/startupers/business/data/order"
	database "github.com/dmanias/startupers/business/sys/database/pgx"
	"github.com/dmanias/startupers/business/sys/database/pgx/dbarray"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, mdr moderator.Moderator) error {
	const q = `
	INSERT INTO moderators
		(id, name, instruction, date_created, date_updated)
	VALUES
		(:id, :name, :instruction, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBModerator(mdr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, mdr moderator.Moderator) error {
	const q = `
	UPDATE
		moderators
	SET 
		"name" = :name,
		"instruction" = :instruction,
		"date_updated" = :date_updated
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBModerator(mdr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, mdr moderator.Moderator) error {
	data := struct {
		ModeratorID string `db:"id"`
	}{
		ModeratorID: mdr.ID.String(),
	}

	const q = `
	DELETE FROM
		moderators
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (s *Store) Query(ctx context.Context, filter moderator.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]moderator.Moderator, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		moderators`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbMdrs []dbModerator
	if err := database.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbMdrs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreModeratorSlice(dbMdrs), nil
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter moderator.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		moderators`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := database.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, moderatorID uuid.UUID) (moderator.Moderator, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: moderatorID.String(),
	}

	const q = `
	SELECT
		*
	FROM
		moderators
	WHERE 
		id = :id`

	var dbMdr dbModerator
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbMdr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return moderator.Moderator{}, fmt.Errorf("namedquerystruct: %w", moderator.ErrNotFound)
		}
		return moderator.Moderator{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreModerator(dbMdr), nil
}

// QueryByIDs gets the specified users from the database.
func (s *Store) QueryByIDs(ctx context.Context, moderatorIDs []uuid.UUID) ([]moderator.Moderator, error) {
	ids := make([]string, len(moderatorIDs))
	for i, moderatorID := range moderatorIDs {
		ids[i] = moderatorID.String()
	}

	data := struct {
		ModeratorID interface {
			driver.Valuer
			sql.Scanner
		} `db:"id"`
	}{
		ModeratorID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		*
	FROM
		moderators
	WHERE
		id = ANY(:id)`

	var mdrs []dbModerator
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &mdrs); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreModeratorSlice(mdrs), nil
}
