package ideadb

import (
	"fmt"
	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/data/order"
)

var orderByFields = map[string]string{
	idea.OrderByID:          "id",
	idea.OrderByUserID:      "user_id",
	idea.OrderByTitle:       "title",
	idea.OrderByCategory:    "category",
	idea.OrderByPrivacy:     "privacy",
	idea.OrderByStage:       "stage",
	idea.OrderByDateCreated: "date_created",
	idea.OrderByDateUpdated: "date_updated",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}
	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
