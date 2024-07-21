package postgrp

import (
	"fmt"
	"time"

	"github.com/dmanias/startupers/business/core/post"
	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

type AppPost struct {
	ID          string `json:"id"`
	IdeaID      string `json:"ideaID"`
	AuthorID    string `json:"authorID"`
	Content     string `json:"content"`
	OwnerType   string `json:"ownerType"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppPost(post post.Post) AppPost {
	return AppPost{
		ID:          post.ID.String(),
		IdeaID:      post.IdeaID.String(),
		AuthorID:    post.AuthorID.String(),
		Content:     post.Content,
		OwnerType:   post.OwnerType,
		DateCreated: post.DateCreated.Format(time.RFC3339),
		DateUpdated: post.DateUpdated.Format(time.RFC3339),
	}
}

type AppNewPost struct {
	IdeaID        string   `json:"ideaID" validate:"required"`
	AuthorID      string   `json:"authorID" validate:"required"`
	Content       string   `json:"content"`
	OwnerType     string   `json:"ownerType"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	Privacy       string   `json:"privacy"`
	Collaborators []string `json:"collaborators"`
	AvatarURL     string   `json:"avatarURL"`
	Stage         string   `json:"stage"`
	Inspiration   string   `json:"inspiration"`
}

func toCoreNewPost(app AppNewPost) (post.NewPost, error) {
	ideaID, err := uuid.Parse(app.IdeaID)
	if err != nil {
		return post.NewPost{}, fmt.Errorf("parsing ideaID: %w", err)
	}

	authorID, err := uuid.Parse(app.AuthorID)
	if err != nil {
		return post.NewPost{}, fmt.Errorf("parsing authorID: %w", err)
	}

	np := post.NewPost{
		IdeaID:    ideaID,
		AuthorID:  authorID,
		Content:   app.Content,
		OwnerType: app.OwnerType,
	}

	return np, nil
}

func (app AppNewPost) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

type AppUpdatePost struct {
	ID      *string `json:"id"`
	Content *string `json:"content"`
}

func toCoreUpdatePost(app AppUpdatePost) (post.UpdatePost, error) {
	var id uuid.UUID
	if app.ID != nil {
		var err error
		id, err = uuid.Parse(*app.ID)
		if err != nil {
			return post.UpdatePost{}, fmt.Errorf("parsing ID: %w", err)
		}
	}

	up := post.UpdatePost{
		ID:      &id,
		Content: app.Content,
	}

	return up, nil
}

func (app AppUpdatePost) Validate() error {
	if err := validate.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
