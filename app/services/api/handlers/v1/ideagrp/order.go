package ideagrp

import (
	"errors"
	"net/http"

	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/data/order"
	"github.com/dmanias/startupers/business/sys/validate"
)

var orderByFields = map[string]struct{}{
	idea.OrderByID:          {},
	idea.OrderByUserID:      {},
	idea.OrderByTitle:       {},
	idea.OrderByCategory:    {},
	idea.OrderByPrivacy:     {},
	idea.OrderByStage:       {},
	idea.OrderByDateCreated: {},
	idea.OrderByDateUpdated: {},
}

func parseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, idea.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}
