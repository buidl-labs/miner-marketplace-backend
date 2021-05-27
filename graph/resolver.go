package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/buidl-labs/filecoin-chain-indexer/lens"
	"github.com/go-pg/pg/v10"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB      *pg.DB
	LensAPI lens.API
}
