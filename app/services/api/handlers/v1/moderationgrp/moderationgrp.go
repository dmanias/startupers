// Package moderationgrp maintains the group of handlers for user access.
package moderationgrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/dmanias/startupers/business/data/order"
	"net/http"

	v1 "github.com/dmanias/startupers/business/web/v1"
	"github.com/dmanias/startupers/business/web/v1/paging"
	"github.com/dmanias/startupers/foundation/web"
)

type Handlers struct {
	moderator *moderator.Core
}

func New(moderator *moderator.Core) *Handlers {
	return &Handlers{
		moderator: moderator,
	}
}

// Create adds a new moderator to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewModerator
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	nc, err := toCoreNewModerator(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	mdr, err := h.moderator.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, moderator.ErrUniqueName) {
			return v1.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: mdr[%+v]: %w", mdr, err)
	}

	return web.Respond(ctx, w, toAppModeration(mdr), http.StatusCreated)
}

// QueryByNameHandler is a wrapper function that adapts QueryByName to the expected handler signature.
//func (h *Handlers) QueryByNameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//	name := web.Param(r, "name")
//	instruction, err := h.QueryByName(ctx, name)
//	if err != nil {
//		return err
//	}
//	return web.Respond(ctx, w, instruction, http.StatusOK)
//}

type QueryByNameResponse struct {
	Moderator string `json:"moderator"`
	ID        string `json:"id"`
}

// QueryByNameHandler is a wrapper function that adapts QueryByName to the expected handler signature.
func (h *Handlers) QueryByNameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	name := web.Param(r, "name")
	instruction, id, err := h.QueryByName(ctx, name)
	if err != nil {
		return err
	}

	response := QueryByNameResponse{
		Moderator: instruction,
		ID:        id,
	}

	return web.Respond(ctx, w, response, http.StatusOK)
}

// QueryByName searches for a moderator by name and returns the moderator's instruction.
func (h *Handlers) QueryByName(ctx context.Context, name string) (string, string, error) { // Modify this line
	filter := moderator.QueryFilter{Name: &name}
	moderators, err := h.moderator.Query(ctx, filter, order.By{Field: "name"}, 1, 1)
	if err != nil {
		return "", "", fmt.Errorf("querybyname: %w", err) // Modify this line
	}

	if len(moderators) == 0 {
		return "", "", v1.NewRequestError(fmt.Errorf("moderator '%s' not found", name), http.StatusNotFound) // Modify this line
	}

	instruction := moderators[0].Instruction
	id := moderators[0].ID.String() // Add this line

	return instruction, id, nil // Modify this line
}

// Query returns a list of users with paging.
func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := paging.ParseRequest(r)
	if err != nil {
		return err
	}

	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	moderators, err := h.moderator.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppModerator, len(moderators))
	for i, mdr := range moderators {
		items[i] = toAppModeration(mdr)
	}

	total, err := h.moderator.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}
