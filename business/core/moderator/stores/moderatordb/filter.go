package moderatordb

import (
	"bytes"
	"fmt"
	"github.com/dmanias/startupers/business/core/moderator"
	"strings"
)

func (s *Store) applyFilter(filter moderator.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = *filter.ID
		wc = append(wc, "id = :id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}
	if filter.Instruction != nil {
		data["instruction"] = fmt.Sprintf("%%%s%%", *filter.Instruction)
		wc = append(wc, "instruction LIKE :instruction")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
