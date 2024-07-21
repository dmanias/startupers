package ideadb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dmanias/startupers/business/core/idea"
)

func (s *Store) applyFilter(filter idea.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = filter.ID
		wc = append(wc, "id = :id")
	}

	if filter.UserID != nil {
		data["user_id"] = filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if filter.Title != nil {
		data["title"] = fmt.Sprintf("%%%s%%", *filter.Title)
		wc = append(wc, "title LIKE :title")
	}

	if filter.Category != nil {
		data["category"] = *filter.Category
		wc = append(wc, "category = :category")
	}

	if filter.Tag != nil {
		data["tag"] = *filter.Tag
		wc = append(wc, "tags @> ARRAY[:tag]")
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
