package aidb

import (
	"bytes"
	"fmt"
	"github.com/dmanias/startupers/business/core/ai"
	"strings"
)

func (s *Store) applyFilter(filter ai.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = *filter.ID
		wc = append(wc, "id = :id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}
	if filter.Query != nil {
		data["query"] = fmt.Sprintf("%%%s%%", *filter.Query)
		wc = append(wc, "instruction LIKE :query")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
