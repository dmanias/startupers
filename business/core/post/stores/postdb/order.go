package postdb

import (
	"fmt"
	"github.com/dmanias/startupers/business/core/post"
	"github.com/dmanias/startupers/business/data/order"
)

var orderByFields = map[string]string{
	post.OrderByID:          "id",
	post.OrderByIdeaID:      "idea_id",
	post.OrderByAuthorID:    "author_id",
	post.OrderByDateCreated: "date_created",
	post.OrderByDateUpdated: "date_updated",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}
	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
