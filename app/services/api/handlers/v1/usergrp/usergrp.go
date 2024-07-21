// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/business/web/auth"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	"net/mail"
	"time"

	"github.com/dmanias/startupers/business/core/user"
	v1 "github.com/dmanias/startupers/business/web/v1"
	"github.com/dmanias/startupers/business/web/v1/paging"
	"github.com/dmanias/startupers/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user      *user.Core
	auth      *auth.Auth
	ActiveKID string
	log       *zap.SugaredLogger
}

// New constructs a handlers for route access.
func New(user *user.Core, auth *auth.Auth, activeKID string, log *zap.SugaredLogger) *Handlers {
	return &Handlers{
		user:      user,
		auth:      auth,
		ActiveKID: activeKID,
		log:       log,
	}
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUser
	if err := web.Decode(r, &app); err != nil {
		return err
	}
	// Log the received credentials
	h.log.Infow("Create user", "email", app.Email, "password", app.Password)

	nc, err := toCoreNewUser(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return v1.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: user[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
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

	users, err := h.user.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

func (h *Handlers) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := web.Decode(r, &credentials); err != nil {
		return err
	}

	email, err := mail.ParseAddress(credentials.Email)
	if err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}

	usr, err := h.user.Authenticate(ctx, *email, credentials.Password)
	if err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	// Generate JWT token
	// Create a Claims instance
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "BackEnd",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles:    usr.Roles,
		UserName: usr.Name,
	}
	token, err := h.auth.GenerateToken(h.ActiveKID, claims)
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	// Return the token and user's name in the response
	response := map[string]string{
		"token": token,
	}

	return web.Respond(ctx, w, response, http.StatusOK)
}
