package ideadb

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/data/order"
	database "github.com/dmanias/startupers/business/sys/database/pgx"
	"github.com/dmanias/startupers/business/sys/database/pgx/dbarray"
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

func (s *Store) Create(ctx context.Context, idea idea.Idea) error {
	const q = `
	INSERT INTO ideas
		(id, user_id, title, description, category, tags, privacy, collaborators, avatar_url, stage, inspiration, date_created, date_updated)
	VALUES
		(:id, :user_id, :title, :description, :category, :tags, :privacy, :collaborators, :avatar_url, :stage, :inspiration, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBIdea(idea)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}
	fmt.Printf("store idea")
	return nil
}

func (s *Store) Update(ctx context.Context, idea idea.Idea) error {
	const q = `
	UPDATE
		ideas
	SET 
		"title" = :title,
		"description" = :description,
		"category" = :category,
		"tags" = :tags,
		"privacy" = :privacy,
		"collaborators" = :collaborators,
		"avatar_url" = :avatar_url,
		"stage" = :stage,
		"inspiration" = :inspiration,
		"date_updated" = :date_updated
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBIdea(idea)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, idea idea.Idea) error {
	data := struct {
		IdeaID string `db:"id"`
	}{
		IdeaID: idea.ID.String(),
	}

	const q = `
	DELETE FROM
		ideas
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Query(ctx context.Context, filter idea.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]idea.Idea, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		ideas`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbIdeas []dbIdea
	if err := database.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbIdeas); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreIdeaSlice(dbIdeas), nil
}

func (s *Store) Count(ctx context.Context, filter idea.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		ideas`

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

func (s *Store) QueryByID(ctx context.Context, ideaID uuid.UUID) (idea.Idea, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: ideaID.String(),
	}

	const q = `
	SELECT
		*
	FROM
		ideas
	WHERE 
		id = :id`

	var dbIdea dbIdea
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbIdea); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return idea.Idea{}, fmt.Errorf("namedquerystruct: %w", idea.ErrNotFound)
		}
		return idea.Idea{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreIdea(dbIdea), nil
}

func (s *Store) QueryByIDs(ctx context.Context, ideaIDs []uuid.UUID) ([]idea.Idea, error) {
	ids := make([]string, len(ideaIDs))
	for i, ideaID := range ideaIDs {
		ids[i] = ideaID.String()
	}

	data := struct {
		IdeaID interface {
			driver.Valuer
			sql.Scanner
		} `db:"id"`
	}{
		IdeaID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		*
	FROM
		ideas
	WHERE
		id = ANY(:id)`

	var ideas []dbIdea
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &ideas); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, idea.ErrNotFound
		}
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreIdeaSlice(ideas), nil
}

func (s *Store) QueryTags(ctx context.Context) ([]string, error) {
	const q = `
	SELECT DISTINCT unnest(tags) as tag
	FROM ideas`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("querying tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scanning tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return tags, nil
}
