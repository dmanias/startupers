package challengegrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dmanias/startupers/business/core/challenge"
	v1 "github.com/dmanias/startupers/business/web/v1"
	"github.com/dmanias/startupers/business/web/v1/paging"
	"github.com/dmanias/startupers/foundation/web"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handlers manages the set of challenge endpoints.
type Handlers struct {
	challenge *challenge.Core
	log       *zap.SugaredLogger
}

// New constructs a handlers for route access.
func New(challenge *challenge.Core, log *zap.SugaredLogger) *Handlers {
	return &Handlers{
		challenge: challenge,
		log:       log,
	}
}

func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewChallenge
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	nc, err := toCoreNewChallenge(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	newChallenge, err := h.challenge.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: challenge[%+v]: %w", newChallenge, err)
	}

	return web.Respond(ctx, w, toAppChallenge(newChallenge), http.StatusCreated)
}

func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateChallenge
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	uc, err := toCoreUpdateChallenge(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	challenge, err := h.challenge.QueryByID(ctx, *uc.ID)
	if err != nil {
		return fmt.Errorf("query: challengeID[%s]: %w", uc.ID, err)
	}

	updatedChallenge, err := h.challenge.Update(ctx, challenge, uc)
	if err != nil {
		return fmt.Errorf("update: challenge[%+v]: %w", challenge, err)
	}

	return web.Respond(ctx, w, toAppChallenge(updatedChallenge), http.StatusOK)
}

func (h *Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	challengeID, err := uuid.Parse(web.Param(r, "challenge_id"))
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	challenge, err := h.challenge.QueryByID(ctx, challengeID)
	if err != nil {
		return fmt.Errorf("query: challengeID[%s]: %w", challengeID, err)
	}

	if err := h.challenge.Delete(ctx, challenge); err != nil {
		return fmt.Errorf("delete: challenge[%+v]: %w", challenge, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	challengeID := web.Param(r, "challenge_id")

	id, err := uuid.Parse(challengeID)
	if err != nil {
		return v1.NewRequestError(fmt.Errorf("invalid challenge ID: %w", err), http.StatusBadRequest)
	}

	challenge, err := h.challenge.QueryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("query challenge by ID: %w", err)
	}

	appChallenge := toAppChallenge(challenge)

	return web.Respond(ctx, w, appChallenge, http.StatusOK)
}

// Query returns a list of challenges with paging.
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

	challenges, err := h.challenge.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppChallenge, len(challenges))
	for i, challenge := range challenges {
		items[i] = toAppChallenge(challenge)
	}

	total, err := h.challenge.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}
