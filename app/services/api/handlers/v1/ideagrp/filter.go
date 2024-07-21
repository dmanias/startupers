package ideagrp

import (
	"net/http"
	"time"

	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (idea.QueryFilter, error) {
	values := r.URL.Query()
	var filter idea.QueryFilter

	if ideaID := values.Get("id"); ideaID != "" {
		id, err := uuid.Parse(ideaID)
		if err != nil {
			return idea.QueryFilter{}, validate.NewFieldsError("id", err)
		}
		filter.WithIdeaID(id)
	}

	if userID := values.Get("user_id"); userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return idea.QueryFilter{}, validate.NewFieldsError("user_id", err)
		}
		filter.WithUserID(id)
	}

	if title := values.Get("title"); title != "" {
		filter.WithTitle(title)
	}

	if category := values.Get("category"); category != "" {
		filter.WithCategory(category)
	}

	if tag := values.Get("tag"); tag != "" {
		filter.WithTag(tag)
	}

	if createdDate := values.Get("start_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return idea.QueryFilter{}, validate.NewFieldsError("start_created_date", err)
		}
		filter.WithStartDateCreated(t)
	}

	if createdDate := values.Get("end_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return idea.QueryFilter{}, validate.NewFieldsError("end_created_date", err)
		}
		filter.WithEndCreatedDate(t)
	}

	if err := filter.Validate(); err != nil {
		return idea.QueryFilter{}, err
	}

	return filter, nil
}
