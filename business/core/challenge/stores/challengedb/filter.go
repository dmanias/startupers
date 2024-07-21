package challengedb

import (
	"bytes"
	"github.com/dmanias/startupers/business/core/challenge"
	"strings"
)

func (s *Store) applyFilter(filter challenge.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = filter.ID
		wc = append(wc, "id = :id")
	}

	if filter.IdeaID != nil {
		data["idea_id"] = filter.IdeaID
		wc = append(wc, "idea_id = :idea_id")
	}

	if filter.ModeratorID != nil {
		data["moderator_id"] = filter.ModeratorID
		wc = append(wc, "moderator_id = :moderator_id")
	}

	if filter.StartCreatedDate != nil {
		data["start_date_created"] = filter.StartCreatedDate
		wc = append(wc, "date_created >= :start_date_created")
	}

	if filter.EndCreatedDate != nil {
		data["end_date_created"] = filter.EndCreatedDate
		wc = append(wc, "date_created <= :end_date_created")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
