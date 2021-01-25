package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
)

func (r *queryResolver) Miner(ctx context.Context, id string) (*model.Miner, error) {
	// select * from miner_infos where id=id; (get miner from db)
	m := &model.Miner{
		ID: "f099",
	}
	return m, nil
}

func (r *queryResolver) AllMiners(ctx context.Context) ([]*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) StorageDeal(ctx context.Context, id string) (*model.StorageDeal, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AllStorageDeals(ctx context.Context, since *int, till *int) ([]*model.StorageDeal, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Transaction(ctx context.Context, id string) (*model.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AllTransactions(ctx context.Context, since *int, till *int) ([]*model.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
