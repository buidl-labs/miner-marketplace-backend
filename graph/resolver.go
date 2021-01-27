package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"

	"github.com/go-pg/pg/v10"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos []*model.Todo
	DB    *pg.DB
}
