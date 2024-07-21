// Package challengedb contains user related CRUD functionality.
package aidb

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/business/core/ai"
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
func (s *Store) Create(ctx context.Context, mdr ai.Ai) error {
	const q = `
	INSERT INTO ais
		(id, name, query, date_created, date_updated)
	VALUES
		(:id, :name, :query, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBAi(mdr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, mdr ai.Ai) error {
	const q = `
	UPDATE
		ais
	SET 
		"name" = :name,
		"query" = :query,
		"date_updated" = :date_updated
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBAi(mdr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, mdr ai.Ai) error {
	data := struct {
		AiID string `db:"id"`
	}{
		AiID: mdr.ID.String(),
	}

	const q = `
	DELETE FROM
		ais
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (s *Store) Query(ctx context.Context, filter ai.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]ai.Ai, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		ais`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbAis []dbAi
	if err := database.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbAis); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreAiSlice(dbAis), nil
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter ai.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		ais`

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
func (s *Store) QueryByID(ctx context.Context, aiID uuid.UUID) (ai.Ai, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: aiID.String(),
	}

	const q = `
	SELECT
		*
	FROM
		ais
	WHERE 
		id = :id`

	var dbA dbAi
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbA); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ai.Ai{}, fmt.Errorf("namedquerystruct: %w", ai.ErrNotFound)
		}
		return ai.Ai{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreAi(dbA), nil
}

// QueryByIDs gets the specified users from the database.
func (s *Store) QueryByIDs(ctx context.Context, aiIDs []uuid.UUID) ([]ai.Ai, error) {
	ids := make([]string, len(aiIDs))
	for i, aiID := range aiIDs {
		ids[i] = aiID.String()
	}

	data := struct {
		AiID interface {
			driver.Valuer
			sql.Scanner
		} `db:"id"`
	}{
		AiID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		*
	FROM
		ais
	WHERE
		id = ANY(:id)`

	var mdrs []dbAi
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &mdrs); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreAiSlice(mdrs), nil
}
