package challengedb

import (
	"fmt"
	"github.com/dmanias/startupers/business/core/challenge"
	"github.com/dmanias/startupers/business/data/order"
)

var orderByFields = map[string]string{
	challenge.OrderByID:          "id",
	challenge.OrderByIdeaID:      "idea_id",
	challenge.OrderByModeratorID: "moderator_id",
	challenge.OrderByDateCreated: "date_created",
	challenge.OrderByDateUpdated: "date_updated",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}
	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
