package aidb

import (
	"fmt"
	"github.com/dmanias/startupers/business/core/ai"
	"github.com/dmanias/startupers/business/data/order"
)

var orderByFields = map[string]string{
	ai.OrderByID:   "id",
	ai.OrderByName: "name",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
