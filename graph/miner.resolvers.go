package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"reflect"
	"strconv"

	"github.com/buidl-labs/filecoin-chain-indexer/model/blocks"
	"github.com/buidl-labs/filecoin-chain-indexer/model/market"
	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
	"github.com/buidl-labs/filecoin-chain-indexer/model/miner"
	"github.com/buidl-labs/filecoin-chain-indexer/model/power"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	filecoinbig "github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/v4/actors/builtin"
	mineractor "github.com/filecoin-project/specs-actors/v4/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/v4/actors/util/smoothing"
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
		fmt.Println("get maxheight", err)
		panic(err)
	}
	// TOFIX: mf is at some max height (that's present in db).
	// This should be in the selected epoch range.
	err = r.DB.Model(mf).Where("miner_id = ? and height = ?", minerID, maxHeight).Select()
	if err != nil {
		fmt.Println("mffff", err)
		panic(err)
	}
	mi := new(miner.MinerInfo)
	err = r.DB.Model(mi).Where("miner_id = ?", obj.ID).Limit(1).Select()
	if err != nil {
		fmt.Println("fmid", err)
		panic(err)
	}
	fmt.Println("MINERFUND", mf)

	var totalIncome *big.Int
	var totalExpenditure *big.Int

	totalIncome = big.NewInt(0)
	totalExpenditure = big.NewInt(0)
	totalIncome, totalExpenditure = ComputeIncomeExpenditure(r, obj, mi, since, till)

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
		ExitCode:   txn.ExitCode,
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

func (r *minerResolver) StorageDeals(ctx context.Context, obj *model.Miner, since *int, till *int, offset *int, limit *int) ([]*model.StorageDeal, error) {
	var mdps []market.MarketDealProposal

	limitDefault := 100
	offsetZero := 0
	if limit == nil && offset == nil {
		*limit = limitDefault
		*offset = offsetZero
	} else if limit == nil && offset != nil {
		*limit = limitDefault
	} else if limit != nil && offset == nil {
		*offset = offsetZero
	}

	if since != nil {
		if till != nil {
			err := r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				Where("start_epoch >= ? AND start_epoch <= ?", *since, *till).
				Limit(*limit).
				Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				Where("start_epoch >= ?", *since).
				Limit(*limit).
				Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
		}
	} else {
		if till != nil {
			err := r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				Where("start_epoch <= ?", *till).
				Limit(*limit).
				Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				Limit(*limit).
				Offset(*offset).
				Select()
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

func (r *minerResolver) Transactions(ctx context.Context, obj *model.Miner, since *int, till *int, offset *int, limit *int) ([]*model.Transaction, error) {
	var txns []messages.Transaction
	var mdps []market.MarketDealProposal

	limitDefault := 20
	offsetZero := 0
	if limit == nil && offset == nil {
		limit = &limitDefault
		offset = &offsetZero
	} else if limit == nil && offset != nil {
		limit = &limitDefault
	} else if limit != nil && offset == nil {
		offset = &offsetZero
	}

	// query to take union of deals and transactions at a given start_epoch/height
	// select start_epoch, cast(deal_id as text) from market_deal_proposals where start_epoch=629991 union select height, cid from transactions where height=629991;

	// select start_epoch as height, cast(deal_id as text) as id from market_deal_proposals
	// 	where start_epoch>=629989 and start_epoch<=629993 union
	// 	select height, cid as id from transactions
	// 	where height>=629989 and height<=629993 order by height desc limit 10 offset 0;

	// select start_epoch as height, cast(deal_id as text) as id, null as type, null as amount, null as transferred,
	// client_id as sender, provider_id as receiver, null as miner_tip, null as base_fee_burn, null as method_name,
	// null as actor_name, null as exit_code from market_deal_proposals
	// where start_epoch>=629989 and start_epoch<=629993 union
	// select height, cid as id, 1 as type, amount, transferred, sender, receiver, miner_tip, base_fee_burn,
	// method_name, actor_name, exit_code from transactions where height>=629989 and height<=629993
	// order by height desc limit 2000 ;

	var id, amount, transferred, sender, receiver, minerFee, burnFee, methodName, actorName, storagePricePerEpoch string
	var height, endEpoch, exitCode, transactionType int64
	// mytxn:=new(messages.Transaction)

	// result1 := make([]*model.Transaction, len(txns)+len(mdps))
	var result1 []*model.Transaction
	if since != nil {
		if till != nil {
			rows, err1 := r.PQDB.Query(`
				select start_epoch as height, end_epoch, cast(deal_id as text) as id,
				0 as type, null as amount, null as transferred,
				client_id as sender, provider_id as receiver, null as miner_tip,
				null as base_fee_burn, null as method_name, null as actor_name,
				null as exit_code, storage_price_per_epoch from market_deal_proposals
				where start_epoch>=` + fmt.Sprint(*since) + ` and start_epoch<=` + fmt.Sprint(*till) + ` union
				select height, null as end_epoch, cid as id, 1 as type, amount, transferred, sender,
				receiver, miner_tip, base_fee_burn, method_name, actor_name, exit_code,
				null as storage_price_per_epoch from transactions where height>=` + fmt.Sprint(*since) + ` and height<=` + fmt.Sprint(*till) + `
				order by height desc limit ` + fmt.Sprint(*limit) + ` offset ` + fmt.Sprint(*offset) + `;
			`)
			if err1 != nil {
				panic(err1)
			}
			defer rows.Close()
			fmt.Println("mrows", rows)
			for rows.Next() {
				err1 = rows.Scan(
					&height, &endEpoch, &id, &transactionType, &amount, &transferred,
					&sender, &receiver, &minerFee, &burnFee, &methodName,
					&actorName, &exitCode, &storagePricePerEpoch)
				if err1 != nil {
					panic(err1)
				}
				if transactionType == 1 {
					// txn
					amt := amount
					if amount == "0" {
						if transferred != "0" {
							amt = transferred
						}
					}
					label, direction, gas := DeriveTransactionLabels(methodName, actorName)
					result1 = append(result1, &model.Transaction{
						ID:              id,
						Amount:          amt,
						TransactionType: "txn",
						Label:           label,
						Direction:       direction,
						Gas:             gas,
						Sender:          sender,
						Receiver:        receiver,
						Height:          height,
						MinerFee:        minerFee,
						BurnFee:         burnFee,
						MethodName:      methodName,
						ActorName:       actorName,
						ExitCode:        exitCode,
					})
				} else {
					// deal
					amt, _ := CalculateDealPrice(storagePricePerEpoch, height, endEpoch)
					result1 = append(result1, &model.Transaction{
						ID:              fmt.Sprintf("%v", id),
						Amount:          amt,
						TransactionType: "deal",
						Sender:          sender,
						Receiver:        receiver,
						Height:          height,
					})
				}
				// heightInt, _ := strconv.Atoi(height)
				// minPieceSizeInt, _ := strconv.Atoi(minPieceSize)
				// maxPieceSizeInt, _ := strconv.Atoi(maxPieceSize)
				// mis = append(mis, miner.MinerInfo{
				// 	Height:  int64(heightInt),
				// 	MinerID: minerID,
				// 	Address: address,
				// })
			}
			// err1 := r.DB.Model((*messages.Transaction)(nil)).
			// 	ColumnExpr("").
			// 	Where("").
			// 	Union(r.DB.Model(
			// 		(*market.MarketDealProposal)(nil)).
			// 		Where("")).
			// 	Select()
			// if err1 != nil {
			// 	panic(err1)
			// }
			err := r.DB.Model(&txns).
				Where("miner = ?", obj.ID).
				Where("height >= ? AND height <= ?", *since, *till).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
			err = r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				Where("start_epoch >= ? AND start_epoch <= ?", *since, *till).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&txns).
				Where("miner = ?", obj.ID).
				Where("height >= ?", *since).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
			err = r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				Where("start_epoch >= ?", *since).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
		}
	} else {
		if till != nil {
			err := r.DB.Model(&txns).
				Where("miner = ?", obj.ID).
				Where("height <= ?", *till).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
			err = r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				Where("start_epoch <= ?", *till).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&txns).
				Where("miner = ?", obj.ID).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
			err = r.DB.Model(&mdps).
				Where("provider_id = ?", obj.ID).
				// Limit(*limit).
				// Offset(*offset).
				Select()
			if err != nil {
				panic(err)
			}
		}
	}
	// var transactions []*model.Transaction
	// for _, txn := range txns {
	// 	amt := txn.Amount
	// 	if txn.Amount == "0" {
	// 		if txn.Transferred != "0" {
	// 			amt = txn.Transferred
	// 		}
	// 	}
	// 	transactions = append(transactions, &model.Transaction{
	// 		ID:              txn.Cid,
	// 		Amount:          amt,
	// 		TransactionType: "txn",
	// 		Sender:          txn.Sender,
	// 		Receiver:        txn.Receiver,
	// 		Height:          txn.Height,
	// 		ExitCode:        txn.ExitCode,
	// 	})
	// }
	// for _, sd := range mdps {
	// 	amt, _ := CalculateDealPrice(sd.StoragePricePerEpoch, sd.StartEpoch, sd.EndEpoch)
	// 	transactions = append(transactions, &model.Transaction{
	// 		ID:              fmt.Sprintf("%v", sd.DealID),
	// 		Amount:          amt,
	// 		TransactionType: "deal",
	// 		Sender:          sd.ClientID,
	// 		Receiver:        sd.ProviderID,
	// 		Height:          sd.StartEpoch,
	// 	})
	// }

	/**
	 *	Merge market deals and transactions to a single array.
	 */

	result := make([]*model.Transaction, len(txns)+len(mdps))

	fmt.Println("txns", len(txns), "mdps", len(mdps))
	i := 0
	for len(txns) > 0 && len(mdps) > 0 {
		if txns[0].Height < mdps[0].StartEpoch {
			amt := txns[0].Amount
			if txns[0].Amount == "0" {
				if txns[0].Transferred != "0" {
					amt = txns[0].Transferred
				}
			}
			label, direction, gas := DeriveTransactionLabels(txns[0].MethodName, txns[0].ActorName)
			result[i] = &model.Transaction{
				ID:              txns[0].Cid,
				Amount:          amt,
				TransactionType: "txn",
				Label:           label,
				Direction:       direction,
				Gas:             gas,
				Sender:          txns[0].Sender,
				Receiver:        txns[0].Receiver,
				Height:          txns[0].Height,
				MinerFee:        txns[0].MinerTip,
				BurnFee:         txns[0].BaseFeeBurn,
				MethodName:      txns[0].MethodName,
				ActorName:       txns[0].ActorName,
				ExitCode:        txns[0].ExitCode,
			}
			txns = txns[1:]
		} else {
			amt, _ := CalculateDealPrice(mdps[0].StoragePricePerEpoch, mdps[0].StartEpoch, mdps[0].EndEpoch)
			result[i] = &model.Transaction{
				ID:              fmt.Sprintf("%v", mdps[0].DealID),
				Amount:          amt,
				TransactionType: "deal",
				Sender:          mdps[0].ClientID,
				Receiver:        mdps[0].ProviderID,
				Height:          mdps[0].StartEpoch,
			}
			mdps = mdps[1:]
		}
		i++
	}

	for j := 0; j < len(txns); j++ {
		amt := txns[j].Amount
		if txns[j].Amount == "0" {
			if txns[j].Transferred != "0" {
				amt = txns[j].Transferred
			}
		}
		label, direction, gas := DeriveTransactionLabels(txns[j].MethodName, txns[j].ActorName)
		result[i] = &model.Transaction{
			ID:              txns[j].Cid,
			Amount:          amt,
			TransactionType: "txn",
			Label:           label,
			Direction:       direction,
			Gas:             gas,
			Sender:          txns[j].Sender,
			Receiver:        txns[j].Receiver,
			Height:          txns[j].Height,
			MinerFee:        txns[j].MinerTip,
			BurnFee:         txns[j].BaseFeeBurn,
			MethodName:      txns[j].MethodName,
			ActorName:       txns[j].ActorName,
			ExitCode:        txns[j].ExitCode,
		}
		i++
	}

	for j := 0; j < len(mdps); j++ {
		amt, _ := CalculateDealPrice(mdps[j].StoragePricePerEpoch, mdps[j].StartEpoch, mdps[j].EndEpoch)
		result[i] = &model.Transaction{
			ID:              fmt.Sprintf("%v", mdps[j].DealID),
			Amount:          amt,
			TransactionType: "deal",
			Sender:          mdps[j].ClientID,
			Receiver:        mdps[j].ProviderID,
			Height:          mdps[j].StartEpoch,
		}
		i++
	}

	fmt.Println("offset", *offset, "limit", *limit)
	if *limit > len(result) {
		*limit = len(result) - 1
	}
	if *offset < 0 {
		*offset = 0
	}

	return result[*offset:*limit], nil
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

func (r *minerResolver) EstimatedIncome(ctx context.Context, obj *model.Miner, days int) (*model.EstimatedIncome, error) {
	minerID, _ := address.NewFromString(obj.ID)
	powerActorID, _ := address.NewFromString("f04")
	rewardActorID, _ := address.NewFromString("f02")
	ts, _ := r.LensAPI.ChainHead(context.Background())

	daysUntilEligible := big.NewInt(0)
	minQAP := big.NewInt(10995116277760) // 10 TiB
	minerPower, _ := r.LensAPI.StateMinerPower(context.Background(), minerID, ts.Key())
	cmpR := minerPower.MinerPower.QualityAdjPower.Int.Cmp(minQAP)
	if cmpR == -1 {
		// TODO: find daysUntilEligible

		lastMonthTs, _ := r.LensAPI.ChainGetTipSetByHeight(context.Background(), ts.Height()-30*2880, types.EmptyTSK)
		minerPowerLastMonth, _ := r.LensAPI.StateMinerPower(context.Background(), minerID, lastMonthTs.Key())
		dailyPowerGrowthLastMonth := new(big.Int).Div(
			new(big.Int).Sub(
				minerPower.MinerPower.QualityAdjPower.Int,
				minerPowerLastMonth.MinerPower.QualityAdjPower.Int,
			),
			big.NewInt(int64(30)),
		)

		daysUntilEligible = new(big.Int).Div(
			new(big.Int).Sub(minQAP, minerPower.MinerPower.QualityAdjPower.Int),
			dailyPowerGrowthLastMonth,
		)
	}
	PowerActorState, err := r.LensAPI.StateReadState(context.Background(), powerActorID, ts.Key())
	if err != nil {
		panic(err)
	}
	RewardActorState, err := r.LensAPI.StateReadState(context.Background(), rewardActorID, ts.Key())
	if err != nil {
		panic(err)
	}

	// pas, _ := PowerActorState.State.(power3.State)
	// ras, _ := RewardActorState.State.(reward3.State)

	// type ThisEpochQAPowerSmoothed string
	pas, _ := PowerActorState.State.(map[string]interface{})
	ras, _ := RewardActorState.State.(map[string]interface{})
	fmt.Println(reflect.TypeOf(pas["ThisEpochQAPowerSmoothed"]), " ", pas["ThisEpochQAPowerSmoothed"])
	ThisEpochQAPowerSmoothed, _ := pas["ThisEpochQAPowerSmoothed"].(map[string]interface{})
	fmt.Println(reflect.TypeOf(ThisEpochQAPowerSmoothed), ThisEpochQAPowerSmoothed,
		"pe:", ThisEpochQAPowerSmoothed["PositionEstimate"],
		"ve:", ThisEpochQAPowerSmoothed["VelocityEstimate"])

	ThisEpochRewardSmoothed, _ := ras["ThisEpochRewardSmoothed"].(map[string]interface{})
	fmt.Println(reflect.TypeOf(ThisEpochRewardSmoothed), ThisEpochRewardSmoothed,
		"pe:", ThisEpochRewardSmoothed["PositionEstimate"],
		"ve:", ThisEpochRewardSmoothed["VelocityEstimate"])

	a := ThisEpochQAPowerSmoothed["PositionEstimate"].(string)
	ThisEpochQAPowerSmoothedPositionEstimate, _ := new(big.Int).SetString(a, 10)
	fmt.Println("a:", a, "toa:", reflect.TypeOf(a))

	b := ThisEpochQAPowerSmoothed["VelocityEstimate"].(string)
	ThisEpochQAPowerSmoothedVelocityEstimate, _ := new(big.Int).SetString(b, 10)

	c := ThisEpochRewardSmoothed["PositionEstimate"].(string)
	// c := "36266260308195979333" // initial
	ThisEpochRewardSmoothedPositionEstimate, _ := new(big.Int).SetString(c, 10)
	fmt.Println("c:", c, "toc:", reflect.TypeOf(c))

	d := ThisEpochRewardSmoothed["VelocityEstimate"].(string)
	// d := "-109897758509" // initial
	ThisEpochRewardSmoothedVelocityEstimate, _ := new(big.Int).SetString(d, 10)

	fmt.Println("pas", pas, " old ", reflect.TypeOf(PowerActorState.State), " ", PowerActorState.State)
	fmt.Println("ras", ras, " old ", reflect.TypeOf(RewardActorState.State), " ", RewardActorState.State)
	nwqapP := new(big.Int).Div(ThisEpochQAPowerSmoothedPositionEstimate, new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil))
	nwqapV := new(big.Int).Div(ThisEpochQAPowerSmoothedVelocityEstimate, new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil))
	perEpochRewardP := new(big.Int).Div(ThisEpochRewardSmoothedPositionEstimate, new(big.Int).Mul(big.NewInt(2).Exp(big.NewInt(2), big.NewInt(128), nil), big.NewInt(1e18)))
	perEpochRewardV := new(big.Int).Div(ThisEpochRewardSmoothedVelocityEstimate, new(big.Int).Mul(big.NewInt(2).Exp(big.NewInt(2), big.NewInt(128), nil), big.NewInt(1e18)))

	minerProjectedReward := ProjectFutureReward(days, big.NewInt(int64(100000*math.Pow(2, 30))), nwqapP, nwqapV, perEpochRewardP, perEpochRewardV)
	fmt.Println("minerProjectedReward", minerProjectedReward)

	qaPower := minerPower.MinerPower.QualityAdjPower //filecoinbig.NewInt(int64(100000 * math.Pow(2, 30)))
	fmt.Println("minerqaPower", qaPower)
	nrwd := mineractor.ExpectedRewardForPower(smoothing.FilterEstimate{
		PositionEstimate: filecoinbig.NewFromGo(ThisEpochRewardSmoothedPositionEstimate),
		VelocityEstimate: filecoinbig.NewFromGo(ThisEpochRewardSmoothedVelocityEstimate),
	}, smoothing.FilterEstimate{
		PositionEstimate: filecoinbig.NewFromGo(ThisEpochQAPowerSmoothedPositionEstimate),
		VelocityEstimate: filecoinbig.NewFromGo(ThisEpochQAPowerSmoothedVelocityEstimate),
	}, qaPower, builtin.EpochsInDay*abi.ChainEpoch(days))

	// v := 1e-18
	atto := big.NewInt(1e18)
	minerProjectedReward = nrwd.Int.Div(nrwd.Int, atto)
	fmt.Println("nrwd", nrwd)

	// GET currentEpoch: https://filfox.info/api/v1/tipset/recent?count=1
	resp, err := http.Get("https://filfox.info/api/v1/tipset/recent?count=1")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	type FilfoxBlock struct {
		Cid string
	}

	type FilfoxTipset struct {
		Height       int64
		Timestamp    int64
		MessageCount int64
		Blocks       []FilfoxBlock
	}

	var latestTipset []FilfoxTipset

	if err := json.NewDecoder(resp.Body).Decode(&latestTipset); err != nil {
		panic(err)
	}
	fmt.Println("latestTipset: ", latestTipset)

	currentEpoch := latestTipset[0].Height
	var existingMdps []market.MarketDealProposal
	err = r.DB.Model(&existingMdps).
		Where("provider_id = ?", obj.ID).
		Where("start_epoch <= ?", currentEpoch+int64(days*2880)).
		Where("end_epoch >= ?", currentEpoch).
		Select()
	if err != nil {
		panic(err)
	}

	existingDealsEarnings := int64(0)

	for _, mdp := range existingMdps {
		rangeStart := currentEpoch
		rangeEnd := currentEpoch + int64(days*2880)
		if mdp.StartEpoch > currentEpoch {
			rangeStart = mdp.StartEpoch
		}
		if mdp.EndEpoch < currentEpoch+int64(days*2880) {
			rangeEnd = mdp.EndEpoch
		}
		count := rangeEnd - rangeStart

		storagePricePerEpoch, _ := strconv.ParseInt(mdp.StoragePricePerEpoch, 10, 64)
		existingDealsEarnings += count * storagePricePerEpoch
	}

	pricePerEpochSum := int64(0)
	var lastMonthMdps []market.MarketDealProposal
	err = r.DB.Model(&lastMonthMdps).
		Where("provider_id = ?", obj.ID).
		Where("start_epoch >= ?", currentEpoch-int64(30*2880)).
		Select()
	if err != nil {
		panic(err)
	}

	estimatedFutureDealsEarnings := int64(0)
	lastMonthDealsCount := int64(len(lastMonthMdps))
	if lastMonthDealsCount != 0 {
		for _, mdp := range lastMonthMdps {
			storagePricePerEpoch, _ := strconv.ParseInt(mdp.StoragePricePerEpoch, 10, 64)
			pricePerEpochSum += storagePricePerEpoch
		}
		averagePricePerEpochLastMonth := pricePerEpochSum / lastMonthDealsCount
		estimatedFutureDeals := (lastMonthDealsCount / 30) * int64(days)
		estimatedFutureDealsEarnings = estimatedFutureDeals * averagePricePerEpochLastMonth * int64(days*2880)
	}

	ei := &model.EstimatedIncome{
		DealPayments: &model.DealPayments{
			ExistingDeals:        existingDealsEarnings,
			PotentialFutureDeals: estimatedFutureDealsEarnings,
		},
		BlockRewards: &model.BlockRewards{
			BlockRewards:      minerProjectedReward.Int64(),
			DaysUntilEligible: daysUntilEligible.Int64(),
		},
	}
	return ei, nil
}

func (r *minerResolver) EstimatedExpenditure(ctx context.Context, obj *model.Miner, days int) (*model.EstimatedExpenditure, error) {
	ee := &model.EstimatedExpenditure{
		PreCommitExpiryPenalty: 0,
		UndeclaredFaultPenalty: 0,
		DeclaredFaultPenalty:   0,
		OngoingFaultPenalty:    0,
		TerminationPenalty:     0,
		ConsensusFaultPenalty:  0,
	}
	return ee, nil
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
