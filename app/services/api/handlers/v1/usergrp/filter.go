package usergrp

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/dmanias/startupers/business/core/user"
	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (user.QueryFilter, error) {
	values := r.URL.Query()

	var filter user.QueryFilter

	if userID := values.Get("id"); userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return user.QueryFilter{}, validate.NewFieldsError("id", err)
		}
		filter.WithUserID(id)
	}

	if email := values.Get("email"); email != "" {
		addr, err := mail.ParseAddress(email)
		if err != nil {
			return user.QueryFilter{}, validate.NewFieldsError("email", err)
		}
		filter.WithEmail(*addr)
	}

	if createdDate := values.Get("start_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return user.QueryFilter{}, validate.NewFieldsError("start_created_date", err)
		}
		filter.WithStartDateCreated(t)
	}

	if createdDate := values.Get("end_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return user.QueryFilter{}, validate.NewFieldsError("end_created_date", err)
		}
		filter.WithEndCreatedDate(t)
	}

	if name := values.Get("name"); name != "" {
		filter.WithName(name)
	}

	if err := filter.Validate(); err != nil {
		return user.QueryFilter{}, err
	}

	return filter, nil
}
