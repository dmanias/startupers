package moderatordb

import (
	"fmt"
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/dmanias/startupers/business/data/order"
)

var orderByFields = map[string]string{
	moderator.OrderByID:   "id",
	moderator.OrderByName: "name",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
