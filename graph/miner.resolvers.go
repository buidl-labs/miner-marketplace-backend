package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/go-pg/pg/v10/orm"

	"github.com/buidl-labs/filecoin-chain-indexer/model/blocks"
	"github.com/buidl-labs/filecoin-chain-indexer/model/market"
	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
	"github.com/buidl-labs/filecoin-chain-indexer/model/miner"
	"github.com/buidl-labs/filecoin-chain-indexer/model/power"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
)

func (r *contactResolver) Miner(ctx context.Context, obj *model.Contact) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *financeMetricsResolver) Miner(ctx context.Context, obj *model.FinanceMetrics) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *minerResolver) Owner(ctx context.Context, obj *model.Miner) (*model.Owner, error) {
	mi := new(miner.MinerInfo)
	var maxHeight int
	fmt.Println("minerid", obj.ID)
	err := r.DB.Model(mi).ColumnExpr("max(height)").Where("miner_id = ?", obj.ID).Select(&maxHeight)
	if err != nil {
		panic(err)
	}
	var ownerID string
	err = r.DB.Model((*miner.MinerInfo)(nil)).Column("owner_id").Where("miner_id = ? and height = ?", obj.ID, maxHeight).Select(&ownerID)
	if err != nil {
		fmt.Println("myerr", err)
		panic(err)
	}
	o := &model.Owner{
		ID:    ownerID,
		Actor: model.ActorAccount,
	}
	return o, nil
}

func (r *minerResolver) Worker(ctx context.Context, obj *model.Miner) (*model.Worker, error) {
	mi := new(miner.MinerInfo)
	var maxHeight int
	fmt.Println("minerid", obj.ID)
	err := r.DB.Model(mi).ColumnExpr("max(height)").Where("miner_id = ?", obj.ID).Select(&maxHeight)
	if err != nil {
		panic(err)
	}
	var workerID string
	err = r.DB.Model((*miner.MinerInfo)(nil)).Column("worker_id").Where("miner_id = ? and height = ?", obj.ID, maxHeight).Select(&workerID)
	if err != nil {
		panic(err)
	}
	w := &model.Worker{
		ID:    workerID,
		Miner: obj,
		Actor: model.ActorAccount,
	}
	return w, nil
}

func (r *minerResolver) Contact(ctx context.Context, obj *model.Miner) (*model.Contact, error) {
	c := &model.Contact{
		Miner:   obj,
		Email:   "",
		Slack:   "",
		Website: "",
		Twitter: "",
	}
	return c, nil
}

func (r *minerResolver) ServiceDetails(ctx context.Context, obj *model.Miner) (*model.ServiceDetails, error) {
	mi := new(miner.MinerInfo)
	err := r.DB.Model(mi).Where("miner_id = ?", obj.ID).Select()
	if err != nil {
		panic(err)
	}

	sd := &model.ServiceDetails{
		Miner:           obj,
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
	var bhs []blocks.BlockHeader
	var winsum uint64
	err := r.DB.Model(&bhs).ColumnExpr("SUM(win_count) AS wins").Where("miner_id = ?", obj.ID).Select(&winsum)
	if err != nil {
		panic(err)
	}

	cp := new(power.PowerActorClaim)
	var maxHeight int
	err = r.DB.Model(cp).ColumnExpr("max(height)").Where("miner_id = ?", obj.ID).Select(&maxHeight)
	if err != nil {
		panic(err)
	}
	fmt.Println("maxHeight ", maxHeight)
	err = r.DB.Model(cp).Where("miner_id = ? AND height = ?", obj.ID, maxHeight).Select()
	if err != nil {
		panic(err)
	}

	qi := &model.QualityIndicators{
		Miner:           obj,
		WinCount:        winsum,
		RawBytePower:    cp.RawBytePower,
		QualityAdjPower: cp.QualityAdjPower,
	}
	return qi, nil
}

func (r *minerResolver) FinanceMetrics(ctx context.Context, obj *model.Miner, since *int, till *int) (*model.FinanceMetrics, error) {
	minerID := obj.ID
	mf := new(miner.MinerFund)
	var maxHeight int
	fmt.Println("minerid", obj.ID)
	err := r.DB.Model(mf).ColumnExpr("max(height)").Where("miner_id = ?", obj.ID).Select(&maxHeight)
	if err != nil {
		panic(err)
	}
	err = r.DB.Model(mf).Where("miner_id = ? and height = ?", minerID, maxHeight).Select()
	if err != nil {
		panic(err)
	}

	// txn := new(messages.Transaction)
	var txns []messages.Transaction
	var txns1 []messages.Transaction
	var totalIncome *big.Int
	var amts []string
	err = r.DB.Model(&txns).Column("amount").Where("receiver = ?", obj.ID).Select(&amts)
	if err != nil {
		panic(err)
	}
	fmt.Println("amts", amts)
	totalIncome = big.NewInt(0)
	for _, amt := range amts {
		n := new(big.Int)
		n, ok := n.SetString(amt, 10)
		if !ok {
			fmt.Println("SetString: error")
		}
		fmt.Println(n)
		totalIncome.Add(totalIncome, n)
	}
	// err = r.DB.Model(&txns).ColumnExpr("SUM(amount::bigint) AS ti").Where("receiver = ?", obj.ID).Select(&totalIncome)
	// if err != nil {
	// 	panic(err)
	// }
	var totalExpenditure *big.Int
	var amts1 []string
	err = r.DB.Model(&txns1).Column("amount").Where("sender = ?", obj.ID).Select(&amts1)
	if err != nil {
		panic(err)
	}
	fmt.Println("amts1", amts1)
	totalExpenditure = big.NewInt(0)
	for _, amt := range amts1 {
		n := new(big.Int)
		n, ok := n.SetString(amt, 10)
		if !ok {
			fmt.Println("SetString: error")
		}
		fmt.Println(n)
		totalExpenditure.Add(totalExpenditure, n)
	}
	// err = r.DB.Model(&txns1).ColumnExpr("SUM(amount::bigint) AS te").Where("sender = ?", obj.ID).Select(&totalExpenditure)
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println("inc ", totalIncome, " exp ", totalExpenditure)
	fm := &model.FinanceMetrics{
		TotalIncome:           totalIncome.String(),      // strconv.Itoa(int(totalIncome)),
		TotalExpenditure:      totalExpenditure.String(), //  strconv.Itoa(int(totalExpenditure)),
		BlockRewards:          "",
		StorageDealPayments:   "",
		RetrievalDealPayments: "",
		NetworkFee:            "",
		Penalty:               "",
		PreCommitDeposits:     mf.PreCommitDeposits,
		InitialPledge:         mf.InitialPledge,
		LockedFunds:           mf.LockedFunds,
		AvailableFunds:        mf.AvailableBalance,
	}
	return fm, nil
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
	mdp := new(market.MarketDealProposal)
	err := r.DB.Model(mdp).Where("deal_id = ?", id).Select()
	if err != nil {
		panic(err)
	}
	storagedeal := &model.StorageDeal{
		ID:                int(mdp.DealID),
		ClientID:          mdp.ClientID,
		ProviderID:        mdp.ProviderID,
		StartEpoch:        mdp.StartEpoch,
		EndEpoch:          mdp.EndEpoch,
		PaddedPieceSize:   mdp.PaddedPieceSize,
		UnPaddedPieceSize: mdp.UnpaddedPieceSize,
		PieceCid:          mdp.PieceCID,
		Verified:          mdp.IsVerified,
		Miner:             obj,
	}
	return storagedeal, nil
}

func (r *minerResolver) Transaction(ctx context.Context, obj *model.Miner, id string) (*model.Transaction, error) {
	txn := new(messages.Transaction)
	err := r.DB.Model(txn).Where("cid = ?", id).Select()
	if err != nil {
		panic(err)
	}
	transaction := &model.Transaction{
		ID:              txn.Cid,
		Miner:           obj,
		Amount:          txn.Amount,
		Sender:          txn.Sender,
		Receiver:        txn.Receiver,
		Height:          txn.Height,
		NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
		TransactionType: GetTransactionType(txn.MethodName),
	}

	return transaction, nil
}

func (r *minerResolver) Sector(ctx context.Context, obj *model.Miner, id string) (*model.Sector, error) {
	sec := new(miner.MinerSectorInfo)
	fmt.Println("mid ", obj.ID, " sid ", id)
	err := r.DB.Model(sec).Where("sector_id = ? AND miner_id = ?", id, obj.ID).Select()
	if err != nil {
		panic(err)
	}
	return &model.Sector{
		ID:              id,
		Miner:           obj,
		Size:            "",
		ActivationEpoch: sec.ActivationEpoch,
		ExpirationEpoch: sec.ExpirationEpoch,
		InitialPledge:   sec.InitialPledge,
	}, nil
}

func (r *minerResolver) Penalty(ctx context.Context, obj *model.Miner, id string) (*model.Penalty, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) Deadline(ctx context.Context, obj *model.Miner, id string) (*model.Deadline, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *minerResolver) StorageDeals(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.StorageDeal, error) {
	var mdps []market.MarketDealProposal
	if since != nil {
		if till != nil {
			err := r.DB.Model(&mdps).Where("start_epoch >= ? AND start_epoch <= ? AND provider_id = ?", *since, *till, obj.ID).Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&mdps).Where("start_epoch >= ? AND provider_id = ?", *since, obj.ID).Select()
			if err != nil {
				panic(err)
			}
		}
	} else {
		if till != nil {
			err := r.DB.Model(&mdps).Where("start_epoch <= ? AND provider_id = ?", *till, obj.ID).Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&mdps).Where("provider_id = ?", obj.ID).Select()
			if err != nil {
				panic(err)
			}
		}
	}

	var storagedeals []*model.StorageDeal
	for _, mdp := range mdps {
		storagedeals = append(storagedeals, &model.StorageDeal{
			ID:                int(mdp.DealID),
			ClientID:          mdp.ClientID,
			ProviderID:        mdp.ProviderID,
			StartEpoch:        mdp.StartEpoch,
			EndEpoch:          mdp.EndEpoch,
			PaddedPieceSize:   mdp.PaddedPieceSize,
			UnPaddedPieceSize: mdp.UnpaddedPieceSize,
			PieceCid:          mdp.PieceCID,
			Verified:          mdp.IsVerified,
			Miner:             obj,
		})
	}
	return storagedeals, nil
}

func (r *minerResolver) Transactions(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.Transaction, error) {
	mi := new(miner.MinerInfo)
	err := r.DB.Model(mi).Where("miner_id = ?", obj.ID).Select()
	if err != nil {
		panic(err)
	}

	var txns []messages.Transaction
	if since != nil {
		if till != nil {
			err := r.DB.Model(&txns).Where("height >= ? AND height <= ?", *since, *till).
				WhereGroup(func(q *orm.Query) (*orm.Query, error) {
					q = q.WhereOr("sender = ? OR receiver = ?",
						obj.ID, obj.ID).
						WhereOr("sender = ? OR receiver = ?",
							mi.OwnerID, mi.OwnerID).
						WhereOr("sender = ? OR receiver = ?",
							mi.WorkerID, mi.WorkerID)
					return q, nil
				}).Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&txns).Where("height >= ?", *since).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.WhereOr("sender = ? OR receiver = ?",
					obj.ID, obj.ID).
					WhereOr("sender = ? OR receiver = ?",
						mi.OwnerID, mi.OwnerID).
					WhereOr("sender = ? OR receiver = ?",
						mi.WorkerID, mi.WorkerID)
				return q, nil
			}).Select()
			if err != nil {
				panic(err)
			}
		}
	} else {
		if till != nil {
			err := r.DB.Model(&txns).Where("height <= ?", *till).
				WhereGroup(func(q *orm.Query) (*orm.Query, error) {
					q = q.WhereOr("sender = ? OR receiver = ?",
						obj.ID, obj.ID).
						WhereOr("sender = ? OR receiver = ?",
							mi.OwnerID, mi.OwnerID).
						WhereOr("sender = ? OR receiver = ?",
							mi.WorkerID, mi.WorkerID)
					return q, nil
				}).Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&txns).Select()
			if err != nil {
				panic(err)
			}
		}
	}
	// err = r.DB.Model(&txns).Where("sender = ? OR receiver = ?",
	// 	obj.ID, obj.ID).WhereOr("sender = ? OR receiver = ?",
	// 	mi.OwnerID, mi.OwnerID).WhereOr("sender = ? OR receiver = ?",
	// 	mi.WorkerID, mi.WorkerID).Select()
	// if err != nil {
	// 	panic(err)
	// }
	var transactions []*model.Transaction
	for _, txn := range txns {
		transactions = append(transactions, &model.Transaction{
			ID:              txn.Cid,
			Miner:           obj,
			Amount:          txn.Amount,
			Sender:          txn.Sender,
			Receiver:        txn.Receiver,
			Height:          txn.Height,
			NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
			TransactionType: GetTransactionType(txn.MethodName),
		})
	}
	return transactions, nil
}

func (r *minerResolver) Sectors(ctx context.Context, obj *model.Miner, since *int, till *int) ([]*model.Sector, error) {
	var sectors []miner.MinerSectorInfo
	err := r.DB.Model(&sectors).Where("miner_id = ?", obj.ID).Select()
	if err != nil {
		panic(err)
	}
	var sects []*model.Sector
	for _, s := range sectors {
		sects = append(sects, &model.Sector{
			ID:              strconv.Itoa(int(s.SectorID)),
			ActivationEpoch: s.ActivationEpoch,
			ExpirationEpoch: s.ExpirationEpoch,
			Size:            "",
			InitialPledge:   s.InitialPledge,
		})
	}
	return sects, nil
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
	return obj.Miner, nil
}

func (r *qualityIndicatorsResolver) WinCount(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	return int(obj.WinCount), nil
}

func (r *qualityIndicatorsResolver) FaultySectors(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	var msfs []*miner.MinerSectorFault
	var faultyCount int
	err := r.DB.Model(&msfs).ColumnExpr("count(*)").Where("miner_id = ?", obj.Miner.ID).Select(&faultyCount)
	if err != nil {
		panic(err)
	}
	return faultyCount, nil
}

func (r *qualityIndicatorsResolver) BlocksMined(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *qualityIndicatorsResolver) MiningEfficiency(ctx context.Context, obj *model.QualityIndicators) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *sectorResolver) Miner(ctx context.Context, obj *model.Sector) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *sectorResolver) Faults(ctx context.Context, obj *model.Sector) ([]*model.Fault, error) {
	var msfs []*miner.MinerSectorFault
	fmt.Println("mid ", obj.Miner.ID, " sid ", obj.ID)
	err := r.DB.Model(&msfs).Where("sector_id = ? AND miner_id = ?", obj.ID, obj.Miner.ID).Select()
	if err != nil {
		panic(err)
	}
	var faults []*model.Fault
	for _, msf := range msfs {
		faults = append(faults, &model.Fault{
			Height: msf.Height,
		})
	}
	return faults, nil
}

func (r *serviceDetailsResolver) Miner(ctx context.Context, obj *model.ServiceDetails) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *serviceDetailsResolver) MinPieceSize(ctx context.Context, obj *model.ServiceDetails) (int, error) {
	return int(obj.MinPieceSize), nil
}

func (r *serviceDetailsResolver) MaxPieceSize(ctx context.Context, obj *model.ServiceDetails) (int, error) {
	return int(obj.MaxPieceSize), nil
}

func (r *storageDealResolver) Miner(ctx context.Context, obj *model.StorageDeal) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *storageDealResolver) PaddedPieceSize(ctx context.Context, obj *model.StorageDeal) (int, error) {
	return int(obj.PaddedPieceSize), nil
}

func (r *storageDealResolver) UnpaddedPieceSize(ctx context.Context, obj *model.StorageDeal) (int, error) {
	return int(obj.UnPaddedPieceSize), nil
}

func (r *transactionResolver) Miner(ctx context.Context, obj *model.Transaction) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *workerResolver) Miner(ctx context.Context, obj *model.Worker) (*model.Miner, error) {
	return obj.Miner, nil
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
