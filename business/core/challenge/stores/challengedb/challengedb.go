package challengedb

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmanias/startupers/business/core/challenge"
	"github.com/dmanias/startupers/business/data/order"
	database "github.com/dmanias/startupers/business/sys/database/pgx"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s *Store) Create(ctx context.Context, challenge challenge.Challenge) error {
	const q = `
    INSERT INTO challenges
        (id, idea_id, moderator_id, answer, photo_url, date_created, date_updated)
    VALUES
        (:id, :idea_id, :moderator_id, :answer, :photo_url, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBChallenge(challenge)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, challenge challenge.Challenge) error {
	const q = `
    UPDATE
        challenges
    SET
        "answer" = :answer,
        "photo_url" = :photo_url,
        "date_updated" = :date_updated
    WHERE
        id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBChallenge(challenge)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, challenge challenge.Challenge) error {
	data := struct {
		ChallengeID string `db:"id"`
	}{
		ChallengeID: challenge.ID.String(),
	}

	const q = `
    DELETE FROM
        challenges
    WHERE
        id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Query(ctx context.Context, filter challenge.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]challenge.Challenge, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
        SELECT
            *
        FROM
            challenges`

	buf := bytes.NewBufferString(q)
	if filter.IdeaID != nil {
		buf.WriteString(" WHERE idea_id = :idea_id")
		data["idea_id"] = filter.IdeaID.String()
	}

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbChallenges []dbChallenge
	if err := database.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbChallenges); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreChallengeSlice(dbChallenges), nil
}

func (s *Store) Count(ctx context.Context, filter challenge.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
    SELECT
        count(1)
    FROM
        challenges`

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

func (s *Store) QueryByID(ctx context.Context, challengeID uuid.UUID) (challenge.Challenge, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: challengeID.String(),
	}

	const q = `
    SELECT
        *
    FROM
        challenges
    WHERE
        id = :id`

	var dbChallenge dbChallenge
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbChallenge); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return challenge.Challenge{}, challenge.ErrNotFound
		}
		return challenge.Challenge{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreChallenge(dbChallenge), nil
}
