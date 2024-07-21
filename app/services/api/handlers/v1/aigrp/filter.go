package aigrp

import (
	"github.com/dmanias/startupers/business/core/ai"
	"net/http"
	"time"

	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (ai.QueryFilter, error) {
	values := r.URL.Query()

	var filter ai.QueryFilter

	if aiID := values.Get("id"); aiID != "" {
		id, err := uuid.Parse(aiID)
		if err != nil {
			return ai.QueryFilter{}, validate.NewFieldsError("id", err)
		}
		filter.WithAiID(id)
	}

	if name := values.Get("name"); name != "" {
		filter.WithName(name)
	}

	if createdDate := values.Get("start_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return ai.QueryFilter{}, validate.NewFieldsError("start_created_date", err)
		}
		filter.WithStartDateCreated(t)
	}

	if createdDate := values.Get("end_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return ai.QueryFilter{}, validate.NewFieldsError("end_created_date", err)
		}
		filter.WithEndCreatedDate(t)
	}

	if err := filter.Validate(); err != nil {
		return ai.QueryFilter{}, err
	}

	return filter, nil
}
