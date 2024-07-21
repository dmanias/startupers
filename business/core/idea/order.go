package idea

import "github.com/dmanias/startupers/business/data/order"

var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

const (
	OrderByID          = "ideaid"
	OrderByUserID      = "userid"
	OrderByTitle       = "title"
	OrderByCategory    = "category"
	OrderByPrivacy     = "privacy"
	OrderByStage       = "stage"
	OrderByDateCreated = "datecreated"
	OrderByDateUpdated = "dateupdated"
)
