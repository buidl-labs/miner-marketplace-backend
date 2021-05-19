package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"database/sql"

	"github.com/go-pg/pg/v10"
	// pq postgresql driver
	_ "github.com/lib/pq"

	"github.com/buidl-labs/filecoin-chain-indexer/lens"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos   []*model.Todo
	DB      *pg.DB
	PQDB    *sql.DB
	miners  []*model.Miner
	LensAPI lens.API
}
