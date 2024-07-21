package post

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/data/order"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("post not found")
)

type Storer interface {
	Create(ctx context.Context, post Post) error
	Update(ctx context.Context, post Post) error
	Delete(ctx context.Context, post Post) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Post, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, postID uuid.UUID) (Post, error)
}

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) Create(ctx context.Context, np NewPost) (Post, error) {
	now := time.Now()

	post := Post{
		ID:          uuid.New(),
		IdeaID:      np.IdeaID,
		AuthorID:    np.AuthorID,
		Content:     np.Content,
		OwnerType:   np.OwnerType,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, post); err != nil {
		return Post{}, fmt.Errorf("create: %w", err)
	}

	return post, nil
}

func (c *Core) Update(ctx context.Context, post Post, up UpdatePost) (Post, error) {
	if up.Content != nil {
		post.Content = *up.Content
	}
	post.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, post); err != nil {
		return Post{}, fmt.Errorf("update: %w", err)
	}

	return post, nil
}

func (c *Core) Delete(ctx context.Context, post Post) error {
	if err := c.storer.Delete(ctx, post); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Post, error) {
	posts, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return posts, nil
}

func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

func (c *Core) QueryByID(ctx context.Context, postID uuid.UUID) (Post, error) {
	post, err := c.storer.QueryByID(ctx, postID)
	if err != nil {
		return Post{}, fmt.Errorf("query: postID[%s]: %w", postID, err)
	}

	return post, nil
}
