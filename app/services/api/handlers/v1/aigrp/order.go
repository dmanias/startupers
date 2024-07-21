package aigrp

import (
	"errors"
	"github.com/dmanias/startupers/business/core/ai"
	"github.com/dmanias/startupers/business/core/post"
	"net/http"

	"github.com/dmanias/startupers/business/data/order"
	"github.com/dmanias/startupers/business/sys/validate"
)

var orderByFields = map[string]struct{}{
	ai.OrderByID:   {},
	ai.OrderByName: {},
}

func parseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, ai.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}

var postOrderByFields = map[string]struct{}{
	post.OrderByID:          {},
	post.OrderByIdeaID:      {},
	post.OrderByAuthorID:    {},
	post.OrderByDateCreated: {},
	post.OrderByDateUpdated: {},
}

func postParseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, post.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := postOrderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}
