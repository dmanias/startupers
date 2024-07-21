// Package moderator provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/business/data/order"
	"time"

	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueName            = errors.New("name is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, prd Ai) error
	Update(ctx context.Context, prd Ai) error
	Delete(ctx context.Context, prd Ai) error
	Count(ctx context.Context, filter QueryFilter) (int, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Ai, error)
	QueryByID(ctx context.Context, aiID uuid.UUID) (Ai, error)
}

// Core manages the set of APIs for ai access.
type Core struct {
	storer Storer
}

// NewCore constructs a core for ai api access.
func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

// Create adds a Moderator to the database. It returns the created Moderator with
// fields like ID and DateCreated populated.
func (c *Core) Create(ctx context.Context, np NewAi) (Ai, error) {
	now := time.Now()

	prd := Ai{
		ID:          uuid.New(),
		Name:        np.Name,
		Query:       np.Query,
		UserID:      np.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, prd); err != nil {
		return Ai{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies data about a Ai. It will error if the specified ID is
// invalid or does not reference an existing Product.
func (c *Core) Update(ctx context.Context, prd Ai, up UpdateAi) (Ai, error) {
	if up.Name != nil {
		prd.Name = *up.Name
	}
	if up.Query != nil {
		prd.Query = *up.Query
	}
	prd.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, prd); err != nil {
		return Ai{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

// Delete removes the ai identified by a given ID.
func (c *Core) Delete(ctx context.Context, prd Ai) error {
	if err := c.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query gets all Ais from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Ai, error) {
	prds, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// Count returns the total number of ais in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryByID finds the ai identified by a given ID.
func (c *Core) QueryByID(ctx context.Context, aiID uuid.UUID) (Ai, error) {
	prd, err := c.storer.QueryByID(ctx, aiID)
	if err != nil {
		return Ai{}, fmt.Errorf("query: productID[%s]: %w", aiID, err)
	}

	return prd, nil
}
