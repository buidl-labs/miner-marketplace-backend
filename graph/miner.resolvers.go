package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
)

func (r *financeMetricsResolver) Miner(ctx context.Context, obj *model.FinanceMetrics) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *financeMetricsResolver) Income(ctx context.Context, obj *model.FinanceMetrics) (*model.Income, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *financeMetricsResolver) Expenditure(ctx context.Context, obj *model.FinanceMetrics) (*model.Expenditure, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *financeMetricsResolver) Funds(ctx context.Context, obj *model.FinanceMetrics) (*model.Funds, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Owner(ctx context.Context, obj *model.Miner) (*model.Owner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Worker(ctx context.Context, obj *model.Miner) (*model.Worker, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Contact(ctx context.Context, obj *model.Miner) (*model.Contact, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) ServiceDetails(ctx context.Context, obj *model.Miner) (*model.ServiceDetails, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) QualityIndicators(ctx context.Context, obj *model.Miner, since *int, till *int) (*model.QualityIndicators, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) FinanceMetrics(ctx context.Context, obj *model.Miner, since *int, till *int) (*model.FinanceMetrics, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) AllServiceDetails(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.ServiceDetails, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) AllQualityIndicators(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.QualityIndicators, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) AllFinanceMetrics(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.FinanceMetrics, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) StorageDeal(ctx context.Context, obj *model.Miner, id string) (*model.StorageDeal, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Transaction(ctx context.Context, obj *model.Miner, id string) (*model.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Sector(ctx context.Context, obj *model.Miner, id string) (*model.Sector, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Penalty(ctx context.Context, obj *model.Miner, id string) (*model.Penalty, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Deadline(ctx context.Context, obj *model.Miner, id string) (*model.Deadline, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) StorageDeals(ctx context.Context, obj *model.Miner) ([]*model.StorageDeal, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Transactions(ctx context.Context, obj *model.Miner) ([]*model.Transaction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Sectors(ctx context.Context, obj *model.Miner) ([]*model.Sector, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Penalties(ctx context.Context, obj *model.Miner) ([]*model.Penalty, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Deadlines(ctx context.Context, obj *model.Miner) ([]*model.Deadline, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *ownerResolver) Miner(ctx context.Context, obj *model.Owner) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *sectorResolver) Miner(ctx context.Context, obj *model.Sector) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *sectorResolver) Faults(ctx context.Context, obj *model.Sector) ([]*model.Fault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *storageDealResolver) Miner(ctx context.Context, obj *model.StorageDeal) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) Miner(ctx context.Context, obj *model.Transaction) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *workerResolver) Miner(ctx context.Context, obj *model.Worker) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

// FinanceMetrics returns generated.FinanceMetricsResolver implementation.
func (r *Resolver) FinanceMetrics() generated.FinanceMetricsResolver {
	return &financeMetricsResolver{r}
}

// Miner returns generated.MinerResolver implementation.
func (r *Resolver) Miner() generated.MinerResolver { return &minerResolver{r} }

// Owner returns generated.OwnerResolver implementation.
func (r *Resolver) Owner() generated.OwnerResolver { return &ownerResolver{r} }

// Sector returns generated.SectorResolver implementation.
func (r *Resolver) Sector() generated.SectorResolver { return &sectorResolver{r} }

// StorageDeal returns generated.StorageDealResolver implementation.
func (r *Resolver) StorageDeal() generated.StorageDealResolver { return &storageDealResolver{r} }

// Transaction returns generated.TransactionResolver implementation.
func (r *Resolver) Transaction() generated.TransactionResolver { return &transactionResolver{r} }

// Worker returns generated.WorkerResolver implementation.
func (r *Resolver) Worker() generated.WorkerResolver { return &workerResolver{r} }

type financeMetricsResolver struct{ *Resolver }
type minerResolver struct{ *Resolver }
type ownerResolver struct{ *Resolver }
type sectorResolver struct{ *Resolver }
type storageDealResolver struct{ *Resolver }
type transactionResolver struct{ *Resolver }
type workerResolver struct{ *Resolver }
