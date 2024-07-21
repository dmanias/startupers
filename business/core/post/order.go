package post

import "github.com/dmanias/startupers/business/data/order"

var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

const (
	OrderByID          = "postid"
	OrderByIdeaID      = "ideaid"
	OrderByAuthorID    = "authorid"
	OrderByDateCreated = "datecreated"
	OrderByDateUpdated = "dateupdated"
)
