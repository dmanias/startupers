package challengegrp

import (
	"net/http"
	"time"

	"github.com/dmanias/startupers/business/core/challenge"
	"github.com/dmanias/startupers/business/sys/validate"
	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (challenge.QueryFilter, error) {
	values := r.URL.Query()
	var filter challenge.QueryFilter

	if challengeID := values.Get("id"); challengeID != "" {
		id, err := uuid.Parse(challengeID)
		if err != nil {
			return challenge.QueryFilter{}, validate.NewFieldsError("id", err)
		}
		filter.WithChallengeID(id)
	}

	if ideaID := values.Get("idea_id"); ideaID != "" {
		id, err := uuid.Parse(ideaID)
		if err != nil {
			return challenge.QueryFilter{}, validate.NewFieldsError("idea_id", err)
		}
		filter.WithIdeaID(id)
	}

	if moderatorID := values.Get("moderator_id"); moderatorID != "" {
		id, err := uuid.Parse(moderatorID)
		if err != nil {
			return challenge.QueryFilter{}, validate.NewFieldsError("moderator_id", err)
		}
		filter.WithModeratorID(id)
	}

	if createdDate := values.Get("start_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return challenge.QueryFilter{}, validate.NewFieldsError("start_created_date", err)
		}
		filter.WithStartDateCreated(t)
	}

	if createdDate := values.Get("end_created_date"); createdDate != "" {
		t, err := time.Parse(time.RFC3339, createdDate)
		if err != nil {
			return challenge.QueryFilter{}, validate.NewFieldsError("end_created_date", err)
		}
		filter.WithEndCreatedDate(t)
	}

	if err := filter.Validate(); err != nil {
		return challenge.QueryFilter{}, err
	}

	return filter, nil
}
