package moderationgrp

import (
	"github.com/dmanias/startupers/business/core/moderator"
	"net/http"
	"time"

	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (moderator.QueryFilter, error) {
	values := r.URL.Query()

	var filter moderator.QueryFilter

	if moderatorID := values.Get("id"); moderatorID != "" {
		id, err := uuid.Parse(moderatorID)
		if err != nil {
			return moderator.QueryFilter{}, validate.NewFieldsError("id", err)
		}
		filter.WithModeratorID(id)
	}

	if name := values.Get("name"); name != "" {
		filter.WithName(name)
	}

	if createdDate := values.Get("start_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return moderator.QueryFilter{}, validate.NewFieldsError("start_created_date", err)
		}
		filter.WithStartDateCreated(t)
	}

	if createdDate := values.Get("end_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return moderator.QueryFilter{}, validate.NewFieldsError("end_created_date", err)
		}
		filter.WithEndCreatedDate(t)
	}

	if err := filter.Validate(); err != nil {
		return moderator.QueryFilter{}, err
	}

	return filter, nil
}
