package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/buidl-labs/filecoin-chain-indexer/model/blocks"
	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
	"github.com/buidl-labs/filecoin-chain-indexer/model/miner"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
)

func (r *contactResolver) Miner(ctx context.Context, obj *model.Contact) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *financeMetricsResolver) ID(ctx context.Context, obj *model.FinanceMetrics) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

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
	var ownerID string
	err := r.DB.Model((*miner.MinerInfo)(nil)).Column("owner_id").Where("miner_id = ?", obj.ID).Select(&ownerID)
	if err != nil {
		panic(err)
	}
	o := &model.Owner{
		ID: ownerID,
	}
	return o, nil
}

func (r *minerResolver) Worker(ctx context.Context, obj *model.Miner) (*model.Worker, error) {
	var workerID string
	err := r.DB.Model((*miner.MinerInfo)(nil)).Column("worker_id").Where("miner_id = ?", obj.ID).Select(&workerID)
	if err != nil {
		panic(err)
	}
	w := &model.Worker{
		ID: workerID,
	}
	return w, nil
}

func (r *minerResolver) Contact(ctx context.Context, obj *model.Miner) (*model.Contact, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) ServiceDetails(ctx context.Context, obj *model.Miner) (*model.ServiceDetails, error) {
	mi := new(miner.MinerInfo)
	err := r.DB.Model(mi).Where("miner_id = ?", obj.ID).Select()
	if err != nil {
		panic(err)
	}

	sd := &model.ServiceDetails{
		Storage:         true,
		Retrieval:       true,
		Repair:          false,
		OnlineDeals:     true,
		OfflineDeals:    false,
		StorageAskPrice: mi.StorageAskPrice,
		MinPieceSize:    mi.MinPieceSize,
		MaxPieceSize:    mi.MaxPieceSize,
	}
	return sd, nil
}

func (r *minerResolver) QualityIndicators(ctx context.Context, obj *model.Miner, since *int, till *int) (*model.QualityIndicators, error) {
	// select miner_id, sum(win_count) from block_headers group by miner_id;

	var bhs []blocks.BlockHeader
	var winsum uint64
	err := r.DB.Model(&bhs).ColumnExpr("SUM(win_count) AS wins").Where("miner_id = ?", obj.ID).Select(&winsum)
	if err != nil {
		panic(err)
	}

	// var res []struct {
	// 	Wins int64
	// }
	// r.DB.Model(&bhs).ColumnExpr("SUM(win_count) AS wins").Group("miner_id").Select(&res)
	// fmt.Println("res", res, &res)
	// panic(fmt.Errorf("not implemented"))

	qi := &model.QualityIndicators{
		WinCount: winsum,
	}
	return qi, nil
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

func (r *minerResolver) StorageDeals(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.StorageDeal, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Transactions(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.Transaction, error) {
	mi := new(miner.MinerInfo)
	err := r.DB.Model(mi).Where("miner_id = ?", obj.ID).Select()
	if err != nil {
		panic(err)
	}

	var txns []messages.Transaction
	// txn := new(messages.Transaction)
	err = r.DB.Model(&txns).Where("from_addr = ? OR to_addr = ?",
		obj.ID, obj.ID).WhereOr("from_addr = ? OR to_addr = ?",
		mi.OwnerID, mi.OwnerID).WhereOr("from_addr = ? OR to_addr = ?",
		mi.WorkerID, mi.WorkerID).Select()
	if err != nil {
		panic(err)
	}
	var transactions []*model.Transaction
	for _, txn := range txns {
		transactions = append(transactions, &model.Transaction{
			ID:              txn.Cid,
			Amount:          txn.Amount,
			Sender:          txn.FromAddr,
			Receiver:        txn.ToAddr,
			Height:          txn.Height,
			NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
			Timestamp:       time.Now(),
			TransactionType: "",
		})
	}
	return transactions, nil
}

func (r *minerResolver) Sectors(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.Sector, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Penalties(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.Penalty, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Deadlines(ctx context.Context, obj *model.Miner) ([]*model.Deadline, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *ownerResolver) Miners(ctx context.Context, obj *model.Owner) ([]*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *qualityIndicatorsResolver) Miner(ctx context.Context, obj *model.QualityIndicators) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *qualityIndicatorsResolver) WinCount(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	fmt.Println("MIDD", obj.WinCount)
	return int(obj.WinCount), nil
	// var bhs []blocks.BlockHeader
	// var winsum int
	// err := r.DB.Model(&bhs).ColumnExpr("SUM(win_count) AS wins").Where("miner_id = ?", obj.Miner.ID).Select(&winsum)
	// if err != nil {
	// 	panic(err)
	// }
	// return &winsum, nil
}

func (r *qualityIndicatorsResolver) FaultySectors(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *qualityIndicatorsResolver) BlocksMined(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *qualityIndicatorsResolver) MiningEfficiency(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *sectorResolver) Miner(ctx context.Context, obj *model.Sector) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *sectorResolver) Faults(ctx context.Context, obj *model.Sector) ([]*model.Fault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *serviceDetailsResolver) Miner(ctx context.Context, obj *model.ServiceDetails) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *serviceDetailsResolver) MinPieceSize(ctx context.Context, obj *model.ServiceDetails) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *serviceDetailsResolver) MaxPieceSize(ctx context.Context, obj *model.ServiceDetails) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *storageDealResolver) Miner(ctx context.Context, obj *model.StorageDeal) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *storageDealResolver) PaddedPieceSize(ctx context.Context, obj *model.StorageDeal) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *storageDealResolver) UnpaddedPieceSize(ctx context.Context, obj *model.StorageDeal) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *transactionResolver) Miner(ctx context.Context, obj *model.Transaction) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *workerResolver) Miner(ctx context.Context, obj *model.Worker) (*model.Miner, error) {
	panic(fmt.Errorf("not implemented"))
}

// Contact returns generated.ContactResolver implementation.
func (r *Resolver) Contact() generated.ContactResolver { return &contactResolver{r} }

// FinanceMetrics returns generated.FinanceMetricsResolver implementation.
func (r *Resolver) FinanceMetrics() generated.FinanceMetricsResolver {
	return &financeMetricsResolver{r}
}

// Miner returns generated.MinerResolver implementation.
func (r *Resolver) Miner() generated.MinerResolver { return &minerResolver{r} }

// Owner returns generated.OwnerResolver implementation.
func (r *Resolver) Owner() generated.OwnerResolver { return &ownerResolver{r} }

// QualityIndicators returns generated.QualityIndicatorsResolver implementation.
func (r *Resolver) QualityIndicators() generated.QualityIndicatorsResolver {
	return &qualityIndicatorsResolver{r}
}

// Sector returns generated.SectorResolver implementation.
func (r *Resolver) Sector() generated.SectorResolver { return &sectorResolver{r} }

// ServiceDetails returns generated.ServiceDetailsResolver implementation.
func (r *Resolver) ServiceDetails() generated.ServiceDetailsResolver {
	return &serviceDetailsResolver{r}
}

// StorageDeal returns generated.StorageDealResolver implementation.
func (r *Resolver) StorageDeal() generated.StorageDealResolver { return &storageDealResolver{r} }

// Transaction returns generated.TransactionResolver implementation.
func (r *Resolver) Transaction() generated.TransactionResolver { return &transactionResolver{r} }

// Worker returns generated.WorkerResolver implementation.
func (r *Resolver) Worker() generated.WorkerResolver { return &workerResolver{r} }

type contactResolver struct{ *Resolver }
type financeMetricsResolver struct{ *Resolver }
type minerResolver struct{ *Resolver }
type ownerResolver struct{ *Resolver }
type qualityIndicatorsResolver struct{ *Resolver }
type sectorResolver struct{ *Resolver }
type serviceDetailsResolver struct{ *Resolver }
type storageDealResolver struct{ *Resolver }
type transactionResolver struct{ *Resolver }
type workerResolver struct{ *Resolver }
