// Package aigrp maintains the group of handlers for user access.
package aigrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/business/core/ai"
	"net/http"

	v1 "github.com/dmanias/startupers/business/web/v1"
	"github.com/dmanias/startupers/business/web/v1/paging"
	"github.com/dmanias/startupers/foundation/web"
)

type Handlers struct {
	ai *ai.Core
}

func New(ai *ai.Core) *Handlers {
	return &Handlers{
		ai: ai,
	}
}

type Question struct {
	Query string `json:"query"`
}

func HandleAskWrapper(apiKey string) func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := handleAsk(apiKey, ctx, w, r)
		if err != nil {
			// Log the error, send a response to the client, etc.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		return nil
	}
}

func handleAsk(apiKey string, ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var question Question
	if err := web.Decode(r, &question); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}

	aiResponse, err := AskAI(apiKey, question.Query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	return web.Respond(ctx, w, aiResponse, http.StatusOK)

	//response, err := json.Marshal(aiResponse)
	//if err != nil {
	//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	return err
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//_, err = w.Write(response)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewAi
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	nc, err := toCoreNewAi(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	a, err := h.ai.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, ai.ErrUniqueName) {
			return v1.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", a, err)
	}

	return web.Respond(ctx, w, toAppAi(a), http.StatusCreated)
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

	ais, err := h.ai.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppAi, len(ais))
	for i, a := range ais {
		items[i] = toAppAi(a)
	}

	total, err := h.ai.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}
