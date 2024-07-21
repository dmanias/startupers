// Package moderator provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package moderator

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
	ErrNotFound   = errors.New("user not found")
	ErrUniqueName = errors.New("email is not unique")
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, prd Moderator) error
	Update(ctx context.Context, prd Moderator) error
	Delete(ctx context.Context, prd Moderator) error
	Count(ctx context.Context, filter QueryFilter) (int, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Moderator, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Moderator, error)
}

// Core manages the set of APIs for product access.
type Core struct {
	storer Storer
}

// NewCore constructs a core for user api access.
func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

// Create adds a Moderator to the database. It returns the created Moderator with
// fields like ID and DateCreated populated.
func (c *Core) Create(ctx context.Context, np NewModerator) (Moderator, error) {
	now := time.Now()

	prd := Moderator{
		ID:          uuid.New(),
		Name:        np.Name,
		Instruction: np.Instruction,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, prd); err != nil {
		return Moderator{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func (c *Core) Update(ctx context.Context, prd Moderator, up UpdateModerator) (Moderator, error) {
	if up.Name != nil {
		prd.Name = *up.Name
	}
	if up.Instruction != nil {
		prd.Instruction = *up.Instruction
	}
	prd.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, prd); err != nil {
		return Moderator{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

// Delete removes the product identified by a given ID.
func (c *Core) Delete(ctx context.Context, prd Moderator) error {
	if err := c.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query gets all Products from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Moderator, error) {
	prds, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// QueryByID finds the product identified by a given ID.
func (c *Core) QueryByID(ctx context.Context, productID uuid.UUID) (Moderator, error) {
	prd, err := c.storer.QueryByID(ctx, productID)
	if err != nil {
		return Moderator{}, fmt.Errorf("query: productID[%s]: %w", productID, err)
	}

	return prd, nil
}

// Count returns the total number of ais in the store.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}
