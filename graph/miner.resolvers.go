package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/buidl-labs/filecoin-chain-indexer/model/blocks"
	"github.com/buidl-labs/filecoin-chain-indexer/model/market"
	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
	"github.com/buidl-labs/filecoin-chain-indexer/model/miner"
	"github.com/buidl-labs/filecoin-chain-indexer/model/power"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
	"github.com/go-pg/pg/v10/orm"
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
	err := r.DB.Model(mi).Where("miner_id = ?", obj.ID).Limit(1).Select()
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
	var bhs2 []blocks.BlockHeader
	var blocksMined uint64
	err = r.DB.Model(&bhs2).ColumnExpr("COUNT(miner_id) AS blocksmined").Where("miner_id = ?", obj.ID).Select(&blocksMined)
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
	// datastored
	// select sum(padded_piece_size) from market_deal_proposals where provider_id = 'f023534';
	qi := &model.QualityIndicators{
		Miner:           obj,
		WinCount:        winsum,
		BlocksMined:     blocksMined,
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

	var totalIncome *big.Int
	var totalExpenditure *big.Int

	totalIncome = big.NewInt(0)
	totalExpenditure = big.NewInt(0)

	// var txns []messages.Transaction
	// var amts []string
	// var minertips []string
	// var basefeeburns []string
	// var transfers []string
	// err = r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
	// 	Where("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 5, obj.ID).
	// 	WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 3, obj.ID).
	// 	WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 23, obj.ID).
	// 	Select(&amts, &minertips, &basefeeburns, &transfers)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("amts", amts)
	// // totalIncome = big.NewInt(0)
	// totalIncome = ComputeBigIntSum(totalIncome, amts)
	// // err = r.DB.Model(&txns).ColumnExpr("SUM(amount::bigint) AS ti").Where("receiver = ?", obj.ID).Select(&totalIncome)
	// // if err != nil {
	// // 	panic(err)
	// // }

	ownerDidChange := false

	workerDidChange := false

	minerAddressChanges := GetMinerAddressChanges()
	val, ok := minerAddressChanges[obj.ID]
	if ok {
		fmt.Println("hhh", val)
		if len(val.WorkerChanges) != 0 {
			workerDidChange = true
		}
		if len(val.OwnerChanges) != 0 {
			ownerDidChange = true
		}
	} else {
		fmt.Println("hhhcantfind")
	}

	if ownerDidChange {
		var txns []messages.Transaction
		var amts []string
		var minertips []string
		var basefeeburns []string
		var transfers []string

		sort.Slice(val.OwnerChanges[:], func(i, j int) bool {
			return val.OwnerChanges[i].Epoch < val.OwnerChanges[j].Epoch
		})
		oc := val.OwnerChanges[0]

		err = r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
			Where("height <= ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, "fil/3/storagemarket", 3, oc.From).
			WhereOr("height <= ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, "fil/3/storageminer", 16, obj.ID).
			Select(&amts, &minertips, &basefeeburns, &transfers)
		if err != nil {
			panic(err)
		}
		totalIncome = ComputeBigIntSum(totalIncome, amts)
		totalIncome = ComputeBigIntSum(totalIncome, minertips)
		totalIncome = ComputeBigIntSum(totalIncome, basefeeburns)
		totalIncome = ComputeBigIntSum(totalIncome, transfers)

		for i, oc := range val.OwnerChanges {
			if i != len(val.OwnerChanges)-1 {
				// if not last
				err := r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, val.OwnerChanges[i+1], "fil/3/storagemarket", 3, oc.From).
					WhereOr("height >= ? AND height < ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, val.OwnerChanges[i+1], "fil/3/storageminer", 16, obj.ID).
					Select(&amts, &minertips, &basefeeburns, &transfers)
				if err != nil {
					panic(err)
				}
				totalIncome = ComputeBigIntSum(totalIncome, amts)
				totalIncome = ComputeBigIntSum(totalIncome, minertips)
				totalIncome = ComputeBigIntSum(totalIncome, basefeeburns)
				totalIncome = ComputeBigIntSum(totalIncome, transfers)
			} else {
				err := r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, "fil/3/storagemarket", 3, oc.From).
					WhereOr("height >= ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, "fil/3/storageminer", 16, obj.ID).
					Select(&amts, &minertips, &basefeeburns, &transfers)
				if err != nil {
					panic(err)
				}
				totalIncome = ComputeBigIntSum(totalIncome, amts)
				totalIncome = ComputeBigIntSum(totalIncome, minertips)
				totalIncome = ComputeBigIntSum(totalIncome, basefeeburns)
				totalIncome = ComputeBigIntSum(totalIncome, transfers)
			}
		}
	} else {
		var txns []messages.Transaction
		var amts []string
		var minertips []string
		var basefeeburns []string
		var transfers []string

		err = r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
			Where("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, obj.Owner.ID).
			WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 16, obj.ID).
			Select(&amts, &minertips, &basefeeburns, &transfers)
		if err != nil {
			panic(err)
		}
		totalIncome = ComputeBigIntSum(totalIncome, amts)
		totalIncome = ComputeBigIntSum(totalIncome, minertips)
		totalIncome = ComputeBigIntSum(totalIncome, basefeeburns)
		totalIncome = ComputeBigIntSum(totalIncome, transfers)
	}

	if workerDidChange {
		var txns []messages.Transaction
		var amts []string
		var minertips []string
		var basefeeburns []string
		var transfers []string

		var txns1 []messages.Transaction
		var amts1 []string
		var minertips1 []string
		var basefeeburns1 []string
		var transfers1 []string

		minEpoch := int64(10000000000000)

		sort.Slice(val.WorkerChanges[:], func(i, j int) bool {
			return val.WorkerChanges[i].Epoch < val.WorkerChanges[j].Epoch
		})
		wc := val.WorkerChanges[0]
		err := r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
			Where("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 5).
			WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 3).
			WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 23).
			WhereOr("height <= ? AND sender = ? AND actor_name != ?", wc.Epoch, wc.From, "fil/3/storageminer").
			Select(&amts, &minertips, &basefeeburns, &transfers)
		if err != nil {
			panic(err)
		}
		totalExpenditure = ComputeBigIntSum(totalExpenditure, amts)
		totalExpenditure = ComputeBigIntSum(totalExpenditure, minertips)
		totalExpenditure = ComputeBigIntSum(totalExpenditure, basefeeburns)
		totalExpenditure = ComputeBigIntSum(totalExpenditure, transfers)

		err = r.DB.Model(&txns1).Column("amount", "miner_tip", "base_fee_burn", "transferred").
			Where("height <= ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, "fil/3/storagemarket", 3, wc.From).
			Select(&amts1, &minertips1, &basefeeburns1, &transfers1)
		if err != nil {
			panic(err)
		}
		totalIncome = ComputeBigIntSum(totalIncome, amts1)
		totalIncome = ComputeBigIntSum(totalIncome, minertips1)
		totalIncome = ComputeBigIntSum(totalIncome, basefeeburns1)
		totalIncome = ComputeBigIntSum(totalIncome, transfers1)

		for i, wc := range val.WorkerChanges {
			if wc.Epoch < minEpoch {
				minEpoch = wc.Epoch
			}
			if i != len(val.WorkerChanges)-1 {
				// if not last
				err := r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, val.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 5).
					WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, val.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 3).
					WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, val.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 23).
					WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name != ?", wc.Epoch, val.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer").
					Select(&amts, &minertips, &basefeeburns, &transfers)
				if err != nil {
					panic(err)
				}
				totalExpenditure = ComputeBigIntSum(totalExpenditure, amts)
				totalExpenditure = ComputeBigIntSum(totalExpenditure, minertips)
				totalExpenditure = ComputeBigIntSum(totalExpenditure, basefeeburns)
				totalExpenditure = ComputeBigIntSum(totalExpenditure, transfers)

				err = r.DB.Model(&txns1).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, val.WorkerChanges[i+1].Epoch, "fil/3/storagemarket", 3, wc.To).
					Select(&amts1, &minertips1, &basefeeburns1, &transfers1)
				if err != nil {
					panic(err)
				}
				totalIncome = ComputeBigIntSum(totalIncome, amts1)
				totalIncome = ComputeBigIntSum(totalIncome, minertips1)
				totalIncome = ComputeBigIntSum(totalIncome, basefeeburns1)
				totalIncome = ComputeBigIntSum(totalIncome, transfers1)

			} else {
				err := r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 5).
					WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 3).
					WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 23).
					WhereOr("height >= ? AND sender = ? AND actor_name != ?", wc.Epoch, wc.To, "fil/3/storageminer").
					Select(&amts, &minertips, &basefeeburns, &transfers)
				if err != nil {
					panic(err)
				}
				totalExpenditure = ComputeBigIntSum(totalExpenditure, amts)
				totalExpenditure = ComputeBigIntSum(totalExpenditure, minertips)
				totalExpenditure = ComputeBigIntSum(totalExpenditure, basefeeburns)
				totalExpenditure = ComputeBigIntSum(totalExpenditure, transfers)

				err = r.DB.Model(&txns1).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, "fil/3/storagemarket", 3, wc.To).
					Select(&amts1, &minertips1, &basefeeburns1, &transfers1)
				if err != nil {
					panic(err)
				}
				totalIncome = ComputeBigIntSum(totalIncome, amts1)
				totalIncome = ComputeBigIntSum(totalIncome, minertips1)
				totalIncome = ComputeBigIntSum(totalIncome, basefeeburns1)
				totalIncome = ComputeBigIntSum(totalIncome, transfers1)
			}
		}
	} else {
		// FIXME: get worker addr at epoch range
		// Right now we are taking the saved workerID
		// (this is at some random height, needs to be fixed)
		var txns []messages.Transaction
		var amts []string
		var minertips []string
		var basefeeburns []string
		var transfers []string
		err := r.DB.Model(&txns).Column("amount", "miner_tip", "base_fee_burn", "transferred").
			Where("sender = ? AND actor_name = ? AND method != ?", obj.Worker.ID, "fil/3/storageminer", 5).
			WhereOr("sender = ? AND actor_name = ? AND method != ?", obj.Worker.ID, "fil/3/storageminer", 3).
			WhereOr("sender = ? AND actor_name = ? AND method != ?", obj.Worker.ID, "fil/3/storageminer", 23).
			WhereOr("sender = ? AND actor_name != ?", obj.Worker.ID, "fil/3/storageminer").
			Select(&amts, &minertips, &basefeeburns, &transfers)
		if err != nil {
			panic(err)
		}
		totalExpenditure = ComputeBigIntSum(totalExpenditure, amts)
		totalExpenditure = ComputeBigIntSum(totalExpenditure, minertips)
		totalExpenditure = ComputeBigIntSum(totalExpenditure, basefeeburns)
		totalExpenditure = ComputeBigIntSum(totalExpenditure, transfers)

		var txns1 []messages.Transaction
		var amts1 []string
		var minertips1 []string
		var basefeeburns1 []string
		var transfers1 []string
		err = r.DB.Model(&txns1).Column("amount", "miner_tip", "base_fee_burn", "transferred").
			Where("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, obj.Worker.ID).
			Select(&amts1, &minertips1, &basefeeburns1, &transfers1)
		if err != nil {
			panic(err)
		}
	}

	var txns1 []messages.Transaction
	var amts1 []string
	var minertips1 []string
	var basefeeburns1 []string
	var transfers1 []string
	err = r.DB.Model(&txns1).Column("amount", "miner_tip", "base_fee_burn", "transferred").
		Where("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 5, obj.ID).
		WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 3, obj.ID).
		WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 23, obj.ID).
		Select(&amts1, &minertips1, &basefeeburns1, &transfers1)
	if err != nil {
		panic(err)
	}
	fmt.Println("amts1", amts1)
	// totalExpenditure = big.NewInt(0)
	totalExpenditure = ComputeBigIntSum(totalExpenditure, amts1)
	totalExpenditure = ComputeBigIntSum(totalExpenditure, minertips1)
	totalExpenditure = ComputeBigIntSum(totalExpenditure, basefeeburns1)
	totalExpenditure = ComputeBigIntSum(totalExpenditure, transfers1)

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
	amt := "0"
	if txn.Amount == "0" {
		if txn.Transferred != "0" {
			amt = txn.Transferred
		}
	}
	transaction := &model.Transaction{
		ID:         txn.Cid,
		Miner:      obj,
		Amount:     amt,
		Sender:     txn.Sender,
		Receiver:   txn.Receiver,
		Height:     txn.Height,
		MinerFee:   txn.MinerTip,
		BurnFee:    txn.BaseFeeBurn,
		MethodName: txn.MethodName,
		ActorName:  txn.ActorName,
		// NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
		// TransactionType: GetTransactionType(txn.MethodName),
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
	err := r.DB.Model(mi).Where("miner_id = ?", obj.ID).Limit(1).Select()
	if err != nil {
		panic(err)
	}
	// now get the points where owner/worker/control changed

	ownerDidChange := false
	workerDidChange := false
	controlDidChange := false
	minerAddressChanges := GetMinerAddressChanges()
	currMAC, ok := minerAddressChanges[obj.ID]
	if ok {
		if len(currMAC.ControlChanges) != 0 {
			controlDidChange = true
		}
		if len(currMAC.OwnerChanges) != 0 {
			ownerDidChange = true
		}
		if len(currMAC.WorkerChanges) != 0 {
			workerDidChange = true
		}
	} else {
		// NOTE: Assumption is that the whole chain is indexed.
		fmt.Println("no changes")
	}
	fmt.Println(ownerDidChange, workerDidChange, controlDidChange)

	var incoming []messages.Transaction
	var outgoing []messages.Transaction
	var txns []messages.Transaction
	if since != nil {
		if till != nil {
			// Find incoming txns
			if ownerDidChange && workerDidChange {
				var ownerIncoming []messages.Transaction
				sort.Slice(currMAC.OwnerChanges[:], func(i, j int) bool {
					return currMAC.OwnerChanges[i].Epoch < currMAC.OwnerChanges[j].Epoch
				})
				oc := currMAC.OwnerChanges[0]
				var ownerIncoming0 []messages.Transaction

				err = r.DB.Model(&ownerIncoming0).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					WhereGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.
							WhereOr("height <= ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, "fil/3/storagemarket", 3, oc.From).
							WhereOr("height <= ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, "fil/3/storageminer", 16, obj.ID)
						return q, nil
					}).Select()
				if err != nil {
					panic(err)
				}
				ownerIncoming = append(ownerIncoming, ownerIncoming0...)
				for i, oc := range currMAC.OwnerChanges {
					if i != len(currMAC.OwnerChanges)-1 {
						var ownerIncomingNotLast []messages.Transaction
						err := r.DB.Model(&ownerIncomingNotLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							WhereGroup(func(q *orm.Query) (*orm.Query, error) {
								q = q.
									WhereOr("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, currMAC.OwnerChanges[i+1], "fil/3/storagemarket", 3, oc.From).
									WhereOr("height >= ? AND height < ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, currMAC.OwnerChanges[i+1], "fil/3/storageminer", 16, obj.ID)
								return q, nil
							}).Select()
						if err != nil {
							panic(err)
						}
						ownerIncoming = append(ownerIncoming, ownerIncomingNotLast...)
					} else {
						var ownerIncomingLast []messages.Transaction
						err := r.DB.Model(&ownerIncomingLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							WhereGroup(func(q *orm.Query) (*orm.Query, error) {
								q = q.
									WhereOr("height >= ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, "fil/3/storagemarket", 3, oc.From).
									WhereOr("height >= ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, "fil/3/storageminer", 16, obj.ID)
								return q, nil
							}).Select()
						if err != nil {
							panic(err)
						}
						ownerIncoming = append(ownerIncoming, ownerIncomingLast...)
					}
				}

				incoming = append(incoming, ownerIncoming...)

				var workerIncoming []messages.Transaction
				sort.Slice(currMAC.WorkerChanges[:], func(i, j int) bool {
					return currMAC.WorkerChanges[i].Epoch < currMAC.WorkerChanges[j].Epoch
				})
				wc := currMAC.WorkerChanges[0]
				var workerIncoming0 []messages.Transaction

				err = r.DB.Model(&workerIncoming0).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					Where("height <= ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, "fil/3/storagemarket", 3, wc.From).
					Select()
				if err != nil {
					panic(err)
				}
				workerIncoming = append(workerIncoming, workerIncoming0...)
				for i, wc := range currMAC.WorkerChanges {
					if i != len(currMAC.WorkerChanges)-1 {
						var workerIncomingNotLast []messages.Transaction
						err := r.DB.Model(&workerIncomingNotLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, currMAC.WorkerChanges[i+1], "fil/3/storagemarket", 3, wc.From).
							Select()
						if err != nil {
							panic(err)
						}
						workerIncoming = append(workerIncoming, workerIncomingNotLast...)
					} else {
						var workerIncomingLast []messages.Transaction
						err := r.DB.Model(&workerIncomingLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, "fil/3/storagemarket", 3, wc.From).
							Select()
						if err != nil {
							panic(err)
						}
						workerIncoming = append(workerIncoming, workerIncomingLast...)
					}
				}

				incoming = append(incoming, workerIncoming...)
			} else if ownerDidChange && !workerDidChange {
				var ownerIncoming []messages.Transaction
				sort.Slice(currMAC.OwnerChanges[:], func(i, j int) bool {
					return currMAC.OwnerChanges[i].Epoch < currMAC.OwnerChanges[j].Epoch
				})
				oc := currMAC.OwnerChanges[0]
				var ownerIncoming0 []messages.Transaction

				err = r.DB.Model(&ownerIncoming0).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					WhereGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.
							WhereOr("height <= ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, "fil/3/storagemarket", 3, oc.From).
							WhereOr("height <= ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, "fil/3/storageminer", 16, obj.ID)
						return q, nil
					}).Select()
				if err != nil {
					panic(err)
				}
				ownerIncoming = append(ownerIncoming, ownerIncoming0...)
				for i, oc := range currMAC.OwnerChanges {
					if i != len(currMAC.OwnerChanges)-1 {
						var ownerIncomingNotLast []messages.Transaction
						err := r.DB.Model(&ownerIncomingNotLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, currMAC.OwnerChanges[i+1], "fil/3/storagemarket", 3, oc.From).
							WhereOr("height >= ? AND height < ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, currMAC.OwnerChanges[i+1], "fil/3/storageminer", 16, obj.ID).
							Select()
						if err != nil {
							panic(err)
						}
						ownerIncoming = append(ownerIncoming, ownerIncomingNotLast...)
					} else {
						var ownerIncomingLast []messages.Transaction
						err := r.DB.Model(&ownerIncomingLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND actor_name = ? AND method = ? AND sender = ?", oc.Epoch, "fil/3/storagemarket", 3, oc.From).
							WhereOr("height >= ? AND actor_name = ? AND method = ? AND receiver = ?", oc.Epoch, "fil/3/storageminer", 16, obj.ID).
							Select()
						if err != nil {
							panic(err)
						}
						ownerIncoming = append(ownerIncoming, ownerIncomingLast...)
					}
				}

				incoming = append(incoming, ownerIncoming...)

				var workerIncoming []messages.Transaction
				err = r.DB.Model(&workerIncoming).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					Where("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, obj.Worker.ID).
					Select()
				if err != nil {
					panic(err)
				}
				incoming = append(incoming, workerIncoming...)
			} else if workerDidChange && !ownerDidChange {
				var workerIncoming []messages.Transaction
				sort.Slice(currMAC.WorkerChanges[:], func(i, j int) bool {
					return currMAC.WorkerChanges[i].Epoch < currMAC.WorkerChanges[j].Epoch
				})
				wc := currMAC.WorkerChanges[0]
				var workerIncoming0 []messages.Transaction

				err = r.DB.Model(&workerIncoming0).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					Where("height <= ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, "fil/3/storagemarket", 3, wc.From).
					Select()
				if err != nil {
					panic(err)
				}
				workerIncoming = append(workerIncoming, workerIncoming0...)
				for i, wc := range currMAC.WorkerChanges {
					if i != len(currMAC.WorkerChanges)-1 {
						var workerIncomingNotLast []messages.Transaction
						err := r.DB.Model(&workerIncomingNotLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, currMAC.WorkerChanges[i+1], "fil/3/storagemarket", 3, wc.From).
							Select()
						if err != nil {
							panic(err)
						}
						workerIncoming = append(workerIncoming, workerIncomingNotLast...)
					} else {
						var workerIncomingLast []messages.Transaction
						err := r.DB.Model(&workerIncomingLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, "fil/3/storagemarket", 3, wc.From).
							Select()
						if err != nil {
							panic(err)
						}
						workerIncoming = append(workerIncoming, workerIncomingLast...)
					}
				}

				incoming = append(incoming, workerIncoming...)

				var ownerIncoming []messages.Transaction

				err = r.DB.Model(&ownerIncoming).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					WhereGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.
							WhereOr("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, obj.Owner.ID).
							WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 16, obj.ID)
						return q, nil
					}).Select()
				if err != nil {
					panic(err)
				}
				incoming = append(incoming, ownerIncoming...)
			} else {
				var ownerIncoming []messages.Transaction

				err = r.DB.Model(&ownerIncoming).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					WhereGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.
							WhereOr("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, obj.Owner.ID).
							WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 16, obj.ID)
						return q, nil
					}).Select()
				if err != nil {
					panic(err)
				}
				incoming = append(incoming, ownerIncoming...)

				var workerIncoming []messages.Transaction
				err = r.DB.Model(&workerIncoming).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					Where("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, obj.Worker.ID).
					Select()
				if err != nil {
					panic(err)
				}
				incoming = append(incoming, workerIncoming...)
			}
			// end {Find incoming txns}

			// Find outgoing txns
			var outgoing0 []messages.Transaction
			err = r.DB.Model(&outgoing0).Column("amount", "miner_tip", "base_fee_burn", "transferred").
				Where("height >= ? AND height <= ?", *since, *till).
				WhereGroup(func(q *orm.Query) (*orm.Query, error) {
					q = q.
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 5, obj.ID).
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 3, obj.ID).
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 23, obj.ID)
					return q, nil
				}).
				Select()
			if err != nil {
				panic(err)
			}
			outgoing = append(outgoing, outgoing0...)

			if workerDidChange {
				var workerOutgoing []messages.Transaction
				sort.Slice(currMAC.WorkerChanges[:], func(i, j int) bool {
					return currMAC.WorkerChanges[i].Epoch < currMAC.WorkerChanges[j].Epoch
				})
				wc := currMAC.WorkerChanges[0]
				var workerOutgoing0 []messages.Transaction

				err := r.DB.Model(&workerOutgoing0).Column("amount", "miner_tip", "base_fee_burn", "transferred").
					Where("height >= ? AND height <= ?", *since, *till).
					WhereGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.
							WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 5).
							WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 3).
							WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 23).
							WhereOr("height <= ? AND sender = ? AND actor_name != ?", wc.Epoch, wc.From, "fil/3/storageminer")

						return q, nil
					}).Select()
				if err != nil {
					panic(err)
				}
				workerOutgoing = append(workerOutgoing, workerOutgoing0...)

				for i, wc := range currMAC.WorkerChanges {
					if i != len(currMAC.WorkerChanges)-1 {
						var workerOutgoingNotLast []messages.Transaction
						err := r.DB.Model(&workerOutgoingNotLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							WhereGroup(func(q *orm.Query) (*orm.Query, error) {
								q = q.
									WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 5).
									WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 3).
									WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 23).
									WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer")
								return q, nil
							}).Select()
						if err != nil {
							panic(err)
						}
						workerOutgoing = append(workerOutgoing, workerOutgoingNotLast...)
					} else {
						var workerOutgoingLast []messages.Transaction
						err := r.DB.Model(&workerOutgoingLast).Column("amount", "miner_tip", "base_fee_burn", "transferred").
							Where("height >= ? AND height <= ?", *since, *till).
							WhereGroup(func(q *orm.Query) (*orm.Query, error) {
								q = q.
									WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 5).
									WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 3).
									WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 23).
									WhereOr("height >= ? AND sender = ? AND actor_name != ?", wc.Epoch, wc.To, "fil/3/storageminer")
								return q, nil
							}).Select()
						if err != nil {
							panic(err)
						}
						workerOutgoing = append(workerOutgoing, workerOutgoingLast...)
					}
				}

				outgoing = append(outgoing, workerOutgoing...)
			}
			// end {Find outgoing txns}
		} else {
			// FIXME: the below logic needs to be changed
			// wrt changing miner addresses.
			err := r.DB.Model(&txns).Where("height >= ?", *since).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.
					WhereOr("sender = ? OR receiver = ?",
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
			// FIXME: the below logic needs to be changed
			// wrt changing miner addresses.
			err := r.DB.Model(&txns).Where("height <= ?", *till).
				WhereGroup(func(q *orm.Query) (*orm.Query, error) {
					q = q.
						WhereOr("sender = ? OR receiver = ?",
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
			// FIXME: the below logic needs to be changed
			// wrt changing miner addresses.
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
	for _, txn := range incoming {
		amt := "0"
		if txn.Amount == "0" {
			if txn.Transferred != "0" {
				amt = txn.Transferred
			}
		}
		transactions = append(transactions, &model.Transaction{
			ID:         txn.Cid,
			Miner:      obj,
			Amount:     amt,
			Sender:     txn.Sender,
			Receiver:   txn.Receiver,
			Height:     txn.Height,
			MinerFee:   txn.MinerTip,
			BurnFee:    txn.BaseFeeBurn,
			MethodName: txn.MethodName,
			ActorName:  txn.ActorName,
			Direction:  "INCOMING",
			// NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
			// TransactionType: GetTransactionType(txn.MethodName),
		})
	}
	for _, txn := range outgoing {
		amt := "0"
		if txn.Amount == "0" {
			if txn.Transferred != "0" {
				amt = txn.Transferred
			}
		}
		transactions = append(transactions, &model.Transaction{
			ID:         txn.Cid,
			Miner:      obj,
			Amount:     amt,
			Sender:     txn.Sender,
			Receiver:   txn.Receiver,
			Height:     txn.Height,
			MinerFee:   txn.MinerTip,
			BurnFee:    txn.BaseFeeBurn,
			MethodName: txn.MethodName,
			ActorName:  txn.ActorName,
			Direction:  "OUTGOING",
			// NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
			// TransactionType: GetTransactionType(txn.MethodName),
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
	return int(obj.BlocksMined), nil
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
