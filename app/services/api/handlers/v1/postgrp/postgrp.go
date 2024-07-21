package postgrp

import (
	"context"
	"fmt"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/aigrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/moderationgrp"
	"github.com/dmanias/startupers/business/core/post"
	v1 "github.com/dmanias/startupers/business/web/v1"
	"github.com/dmanias/startupers/business/web/v1/paging"
	"github.com/dmanias/startupers/foundation/web"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Handlers manages the set of post endpoints.
type Handlers struct {
	post               *post.Core
	log                *zap.SugaredLogger
	aiHandlers         *aigrp.Handlers
	moderationHandlers *moderationgrp.Handlers
}

// New constructs a handlers for route access.
func New(post *post.Core, log *zap.SugaredLogger, aiHandlers *aigrp.Handlers, moderationHandlers *moderationgrp.Handlers) *Handlers {
	return &Handlers{
		post:               post,
		log:                log,
		aiHandlers:         aiHandlers,
		moderationHandlers: moderationHandlers,
	}
}

func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewPost
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	// Check if the ownerType is 'idea'
	if app.OwnerType == "idea" {

		// Retrieve the idea ID from the app.IdeaID field
		ideaID, err := uuid.Parse(app.IdeaID)
		if err != nil {
			return v1.NewRequestError(err, http.StatusBadRequest)
		}
		// Create a filter to query posts that belong to the specific idea.
		filter := post.QueryFilter{
			IdeaID: &ideaID,
		}

		page, err := paging.ParseRequest(r)
		if err != nil {
			return err
		}

		// Parse the order parameter from the request
		orderBy, err := parseOrder(r)
		if err != nil {
			return err
		}

		posts, err := h.post.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
		if err != nil {
			return fmt.Errorf("query: %w", err)
		}

		items := make([]AppPost, len(posts))
		for i, post := range posts {
			items[i] = toAppPost(post)
		}

		// Get the instruction from the moderator
		instruction, _, err := h.moderationHandlers.QueryByName(ctx, "idea-response")
		if err != nil {
			return err
		}

		// Trim leading and trailing white spaces from the content
		app.Content = strings.TrimSpace(app.Content)

		// Construct the question for the AI
		aiQuestion := instruction + " Idea title: " + app.Title + ", " + "Idea description: " + app.Description + ", " + "Idea tags: " + fmt.Sprint(app.Tags)

		// Include the optional fields in the question
		if app.Inspiration != "" {
			aiQuestion += ", Inspiration: " + app.Inspiration
		}
		if app.Stage != "" {
			aiQuestion += ", Stage: " + app.Stage
		}

		postContent := ""
		if len(posts) > 0 {
			// Include the existing posts in the question
			postContent += "Old Posts: "
			for i, post := range posts {
				postContent += fmt.Sprintf("Post%d: %s, ", i+1, post.Content)
			}
			aiQuestion += ", Old Posts- " + postContent
		}

		// Check if the content is empty
		if app.Content != "" {
			aiQuestion += ", User question: " + app.Content
		}
		fmt.Printf("aiQuestion: %s\n", aiQuestion)
		// Call the AI handler to get the response
		aiResponse, err := h.aiHandlers.Gpt(ctx, aiQuestion)
		if err != nil {
			return err
		}

		// Append the AI response to the post content
		aiResponse = strings.Trim(aiResponse, "\"")
		aiResponse = strings.Trim(aiResponse, "'\"")
		app.Content = aiResponse
	}

	nc, err := toCoreNewPost(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	newPost, err := h.post.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: post[%+v]: %w", newPost, err)
	}

	return web.Respond(ctx, w, toAppPost(newPost), http.StatusCreated)
}

func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdatePost
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	uc, err := toCoreUpdatePost(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	post, err := h.post.QueryByID(ctx, *uc.ID)
	if err != nil {
		return fmt.Errorf("query: postID[%s]: %w", uc.ID, err)
	}

	updatedPost, err := h.post.Update(ctx, post, uc)
	if err != nil {
		return fmt.Errorf("update: post[%+v]: %w", post, err)
	}

	return web.Respond(ctx, w, toAppPost(updatedPost), http.StatusOK)
}

func (h *Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	postID, err := uuid.Parse(web.Param(r, "post_id"))
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	post, err := h.post.QueryByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("query: postID[%s]: %w", postID, err)
	}

	if err := h.post.Delete(ctx, post); err != nil {
		return fmt.Errorf("delete: post[%+v]: %w", post, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	postID := web.Param(r, "post_id")

	id, err := uuid.Parse(postID)
	if err != nil {
		return v1.NewRequestError(fmt.Errorf("invalid post ID: %w", err), http.StatusBadRequest)
	}

	post, err := h.post.QueryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("query post by ID: %w", err)
	}

	appPost := toAppPost(post)

	return web.Respond(ctx, w, appPost, http.StatusOK)
}

// Query returns a list of posts with paging.
func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := paging.ParseRequest(r)
	if err != nil {
		return err
	}

	// Extract the idea_id parameter from the URL.
	ideaIDStr := web.Param(r, "idea_id")

	// Parse the idea_id string into a uuid.UUID.
	ideaID, err := uuid.Parse(ideaIDStr)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	// Create a filter to query posts that belong to the specific idea.
	filter := post.QueryFilter{
		IdeaID: &ideaID,
	}

	// Parse the order parameter from the request
	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	posts, err := h.post.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppPost, len(posts))
	for i, post := range posts {
		items[i] = toAppPost(post)
	}

	total, err := h.post.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}
