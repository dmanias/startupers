package idea

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/data/order"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("idea not found")
)

type Storer interface {
	Create(ctx context.Context, idea Idea) error
	Update(ctx context.Context, idea Idea) error
	Delete(ctx context.Context, idea Idea) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Idea, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, ideaID uuid.UUID) (Idea, error)
	QueryByIDs(ctx context.Context, ideaIDs []uuid.UUID) ([]Idea, error)
	QueryTags(ctx context.Context) ([]string, error)
}

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) Create(ctx context.Context, ni NewIdea) (Idea, error) {
	now := time.Now()

	idea := Idea{
		ID:            uuid.New(),
		UserID:        ni.UserID,
		Title:         ni.Title,
		Description:   ni.Description,
		Category:      ni.Category,
		Tags:          ni.Tags,
		Privacy:       ni.Privacy,
		Collaborators: ni.Collaborators,
		AvatarURL:     ni.AvatarURL,
		Stage:         ni.Stage,
		Inspiration:   ni.Inspiration,
		DateCreated:   now,
		DateUpdated:   now,
	}

	if err := c.storer.Create(ctx, idea); err != nil {
		return Idea{}, fmt.Errorf("create: %w", err)
	}

	return idea, nil
}

func (c *Core) Update(ctx context.Context, idea Idea, ui UpdateIdea) (Idea, error) {
	if ui.Title != nil {
		idea.Title = *ui.Title
	}
	if ui.Description != nil {
		idea.Description = *ui.Description
	}
	if ui.Category != nil {
		idea.Category = *ui.Category
	}
	if ui.Tags != nil {
		idea.Tags = ui.Tags
	}
	if ui.Privacy != nil {
		idea.Privacy = *ui.Privacy
	}
	if ui.Collaborators != nil {
		idea.Collaborators = ui.Collaborators
	}
	if ui.AvatarURL != nil {
		idea.AvatarURL = *ui.AvatarURL
	}
	if ui.Stage != nil {
		idea.Stage = *ui.Stage
	}
	if ui.Inspiration != nil {
		idea.Inspiration = *ui.Inspiration
	}
	idea.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, idea); err != nil {
		return Idea{}, fmt.Errorf("update: %w", err)
	}

	return idea, nil
}

func (c *Core) Delete(ctx context.Context, idea Idea) error {
	if err := c.storer.Delete(ctx, idea); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Idea, error) {
	ideas, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return ideas, nil
}

func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

func (c *Core) QueryByID(ctx context.Context, ideaID uuid.UUID) (Idea, error) {
	idea, err := c.storer.QueryByID(ctx, ideaID)
	if err != nil {
		return Idea{}, fmt.Errorf("query: ideaID[%s]: %w", ideaID, err)
	}

	return idea, nil
}

func (c *Core) QueryByIDs(ctx context.Context, ideaIDs []uuid.UUID) ([]Idea, error) {
	ideas, err := c.storer.QueryByIDs(ctx, ideaIDs)
	if err != nil {
		return nil, fmt.Errorf("query: ideaIDs[%s]: %w", ideaIDs, err)
	}

	return ideas, nil
}

func (c *Core) QueryTags(ctx context.Context) ([]string, error) {
	return c.storer.QueryTags(ctx)
}
