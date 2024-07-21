// Package aigrp maintains the group of handlers for user access.
package aigrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/moderationgrp"
	"github.com/dmanias/startupers/business/core/ai"
	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/dmanias/startupers/business/core/post"
	"github.com/dmanias/startupers/business/web/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"

	v1 "github.com/dmanias/startupers/business/web/v1"
	"github.com/dmanias/startupers/business/web/v1/paging"
	"github.com/dmanias/startupers/foundation/web"
)

type Handlers struct {
	ai                 *ai.Core
	cfg                APIMuxConfig
	moderationHandlers *moderationgrp.Handlers
	postCore           *post.Core
	ideaCore           *idea.Core
}

func New(ai *ai.Core, cfg APIMuxConfig, moderationHandlers *moderationgrp.Handlers, postCore *post.Core, ideaCore *idea.Core) *Handlers {
	return &Handlers{
		ai:                 ai,
		cfg:                cfg,
		moderationHandlers: moderationHandlers,
		postCore:           postCore,
		ideaCore:           ideaCore,
	}
}

type Question struct {
	Query string `json:"query"`
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown      chan os.Signal
	Log           *zap.SugaredLogger
	Auth          *auth.Auth
	DB            *sqlx.DB
	APIKey        string
	Build         string
	ModeratorCore *moderator.Core
	AIType        string
}
type AskResponse struct {
	AIResponse string `json:"ai_response"`
}

func (h *Handlers) Ask(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Get the idea ID, question type, and description from the URL parameters
	ideaID := web.Param(r, "idea_id")
	questionType := web.Param(r, "question_type")
	description := web.Param(r, "description")

	if questionType == "" || description == "" || ideaID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return errors.New("question_type, description and ideaID are required")
	}

	var aiResponse string
	var question string

	// Check the question type and construct the appropriate question
	if questionType == "step" {
		page, err := paging.ParseRequest(r)
		if err != nil {
			return err
		}

		// Get the instruction from the moderator
		instruction, _, err := h.moderationHandlers.QueryByName(ctx, "step")
		if err != nil {
			return err
		}
		// Parse the idea ID into a uuid.UUID
		ideaUUID, err := uuid.Parse(ideaID)
		if err != nil {
			return v1.NewRequestError(err, http.StatusBadRequest)
		}

		// Create a filter to query posts that belong to the specific idea
		filter := post.QueryFilter{
			IdeaID: &ideaUUID,
		}
		orderBy, err := postParseOrder(r)
		if err != nil {
			return err
		}

		// Query the posts for the specific idea

		posts, err := h.postCore.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
		if err != nil {
			return fmt.Errorf("query posts: %w", err)
		}

		// Query the idea by ID
		idea, err := h.ideaCore.QueryByID(ctx, ideaUUID)
		if err != nil {
			return fmt.Errorf("query idea by ID: %w", err)
		}

		// Construct the question with the idea details and posts
		question = fmt.Sprintf("%s\nStep's description: %s\nIdea title: %s\nIdea description: %s\nIdea tags: %v\nPosts:\n%s",
			instruction, description, idea.Title, idea.Description, idea.Tags, formatPosts(posts))

		// Call the GPT function to get the AI response
		aiResponse, err = gpt(h.cfg.APIKey, question)
		if err != nil {
			return err
		}
	} else {
		http.Error(w, "Invalid question type", http.StatusBadRequest)
		return errors.New("invalid question type")
	}

	return web.Respond(ctx, w, AskResponse{AIResponse: aiResponse}, http.StatusOK)
}

func formatPosts(posts []post.Post) string {
	var formattedPosts string
	for i, p := range posts {
		formattedPosts += fmt.Sprintf("Post %d: %s\n", i+1, p.Content)
	}
	return formattedPosts
}

// Dalle is a handler that calls the Dalle function from ai.go
func (h *Handlers) Dalle(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Call the Dalle function from ai.go
	return dalle(h.cfg.APIKey, prompt)
}

func (h *Handlers) Gpt(ctx context.Context, prompt string) (string, error) {
	// Call the Dalle function from ai.go
	return gpt(h.cfg.APIKey, prompt)
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
		return fmt.Errorf("create: ai[%+v]: %w", a, err)
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
