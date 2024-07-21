package challenge

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/data/order"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("challenge not found")
)

type Storer interface {
	Create(ctx context.Context, challenge Challenge) error
	Update(ctx context.Context, challenge Challenge) error
	Delete(ctx context.Context, challenge Challenge) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Challenge, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, challengeID uuid.UUID) (Challenge, error)
}

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) Create(ctx context.Context, nc NewChallenge) (Challenge, error) {
	now := time.Now()

	challenge := Challenge{
		ID:          uuid.New(),
		IdeaID:      nc.IdeaID,
		ModeratorID: nc.ModeratorID,
		Answer:      nc.Answer,
		PhotoURL:    nc.PhotoURL,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, challenge); err != nil {
		return Challenge{}, fmt.Errorf("create: %w", err)
	}

	return challenge, nil
}

func (c *Core) Update(ctx context.Context, challenge Challenge, uc UpdateChallenge) (Challenge, error) {
	if uc.Answer != nil {
		challenge.Answer = *uc.Answer
	}
	if uc.PhotoURL != nil {
		challenge.PhotoURL = *uc.PhotoURL
	}
	challenge.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, challenge); err != nil {
		return Challenge{}, fmt.Errorf("update: %w", err)
	}

	return challenge, nil
}

func (c *Core) Delete(ctx context.Context, challenge Challenge) error {
	if err := c.storer.Delete(ctx, challenge); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Challenge, error) {
	challenges, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return challenges, nil
}

func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

func (c *Core) QueryByID(ctx context.Context, challengeID uuid.UUID) (Challenge, error) {
	challenge, err := c.storer.QueryByID(ctx, challengeID)
	if err != nil {
		return Challenge{}, fmt.Errorf("query: challengeID[%s]: %w", challengeID, err)
	}

	return challenge, nil
}
