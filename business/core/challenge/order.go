package challenge

import "github.com/dmanias/startupers/business/data/order"

var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

const (
	OrderByID          = "challengeid"
	OrderByIdeaID      = "ideaid"
	OrderByModeratorID = "moderatorid"
	OrderByDateCreated = "datecreated"
	OrderByDateUpdated = "dateupdated"
)
