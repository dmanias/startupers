package authgrp

import (
	"context"
	"github.com/dmanias/startupers/business/web/auth"
	"github.com/dmanias/startupers/foundation/web"
	"net/http"
)

type Handlers struct {
	Auth *auth.Auth
}

func New(auth *auth.Auth) *Handlers {
	return &Handlers{
		Auth: auth,
	}
}

func (h *Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return token(h.Auth)(ctx, w, r)
}

func (h *Handlers) Authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return authenticate(h.Auth)(ctx, w, r)
}

func (h *Handlers) Authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return authorize(h.Auth)(ctx, w, r)
}

// token is a handler that generates a JWT token.
func token(a *auth.Auth) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		kid := web.Param(r, "kid")
		claims := auth.GetClaims(ctx)

		tkn, err := a.GenerateToken(kid, claims)
		if err != nil {
			return err
		}

		return web.Respond(ctx, w, map[string]string{"token": tkn}, http.StatusOK)
	}
}

// authenticate is a handler that authenticates a user.
func authenticate(a *auth.Auth) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		bearerToken := r.Header.Get("Authorization")

		claims, err := a.Authenticate(ctx, bearerToken)
		if err != nil {
			return err
		}

		return web.Respond(ctx, w, claims, http.StatusOK)
	}
}

// authorize is a handler that authorizes a user.
func authorize(a *auth.Auth) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var input struct {
			UserID string      `json:"userID"`
			Rule   string      `json:"rule"`
			Claims auth.Claims `json:"claims"`
		}
		if err := web.Decode(r, &input); err != nil {
			return err
		}

		if err := a.Authorize(ctx, input.Claims, input.Rule); err != nil {
			return err
		}

		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
}
