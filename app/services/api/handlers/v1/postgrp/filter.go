package postgrp

import (
	"net/http"
	"time"

	"github.com/dmanias/startupers/business/core/post"
	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (post.QueryFilter, error) {
	values := r.URL.Query()
	var filter post.QueryFilter

	if postID := values.Get("id"); postID != "" {
		id, err := uuid.Parse(postID)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError("id", err)
		}
		filter.WithPostID(id)
	}

	if ideaID := values.Get("idea_id"); ideaID != "" {
		id, err := uuid.Parse(ideaID)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError("idea_id", err)
		}
		filter.WithIdeaID(id)
	}

	if authorID := values.Get("author_id"); authorID != "" {
		id, err := uuid.Parse(authorID)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError("author_id", err)
		}
		filter.WithAuthorID(id)
	}

	if ownerType := values.Get("owner_type"); ownerType != "" {
		filter.WithOwnerType(ownerType)
	}

	if createdDate := values.Get("start_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError("start_created_date", err)
		}
		filter.WithStartDateCreated(t)
	}

	if createdDate := values.Get("end_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError("end_created_date", err)
		}
		filter.WithEndCreatedDate(t)
	}

	if err := filter.Validate(); err != nil {
		return post.QueryFilter{}, err
	}

	return filter, nil
}
