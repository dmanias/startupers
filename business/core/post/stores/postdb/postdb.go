package postdb

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmanias/startupers/business/core/post"
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

func (s *Store) Create(ctx context.Context, post post.Post) error {
	const q = `
    INSERT INTO posts
        (id, idea_id, author_id, content, owner_type, date_created, date_updated)
    VALUES
        (:id, :idea_id, :author_id, :content, :owner_type, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBPost(post)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, post post.Post) error {
	const q = `
    UPDATE
        posts
    SET
        "content" = :content,
        "date_updated" = :date_updated
    WHERE
        id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBPost(post)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, post post.Post) error {
	data := struct {
		PostID string `db:"id"`
	}{
		PostID: post.ID.String(),
	}

	const q = `
    DELETE FROM
        posts
    WHERE
        id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Query(ctx context.Context, filter post.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]post.Post, error) {
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
        SELECT
            *
        FROM
            posts`

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

	var dbPosts []dbPost
	if err := database.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPosts); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCorePostSlice(dbPosts), nil
}

func (s *Store) Count(ctx context.Context, filter post.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
    SELECT
        count(1)
    FROM
        posts`

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

func (s *Store) QueryByID(ctx context.Context, postID uuid.UUID) (post.Post, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: postID.String(),
	}

	const q = `
    SELECT
        *
    FROM
        posts
    WHERE
        id = :id`

	var dbPost dbPost
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPost); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return post.Post{}, post.ErrNotFound
		}
		return post.Post{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCorePost(dbPost), nil
}
