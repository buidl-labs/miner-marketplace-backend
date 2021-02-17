package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/buidl-labs/filecoin-chain-indexer/model/market"
	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
	"github.com/buidl-labs/filecoin-chain-indexer/model/miner"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
	// pq postgresql driver
	_ "github.com/lib/pq"
)

func (r *queryResolver) Miner(ctx context.Context, id string) (*model.Miner, error) {
	mi := new(miner.MinerInfo)
	var maxHeight int
	fmt.Println("minerid", id)
	err := r.DB.Model(mi).ColumnExpr("max(height)").Where("miner_id = ?", id).Select(&maxHeight)
	if err != nil {
		panic(err)
	}
	fmt.Println("maxHeight ", maxHeight)
	err = r.DB.Model(mi).Where("height = ? and miner_id = ?", maxHeight, id).Select()
	if err != nil {
		panic(err)
	}

	m := &model.Miner{
		ID:       mi.MinerID,
		Address:  mi.Address,
		PeerID:   mi.PeerID,
		Name:     "",
		Bio:      "",
		Location: "",
		Verified: false,
	}
	return m, nil
}

func (r *queryResolver) AllMiners(ctx context.Context, after *string, first *int, before *string, last *int, since *int, till *int) ([]*model.Miner, error) {
	// TODO: implement *before, *last, *since, *till
	var mis []miner.MinerInfo // this is the indexer model (db)

	var minerID string
	var address string
	var peerID string
	var ownerID string
	var workerID string
	var height string
	var stateRoot string
	var storageAskPrice string
	var minPieceSize string
	var maxPieceSize string

	if after != nil {
		if first != nil {
			rows, err := r.PQDB.Query(`
				with maxht as (
					select miner_id, max(height) maxh
					from miner_infos
					group by miner_id
				)
				select miner_infos.miner_id, miner_infos.address,
				miner_infos.peer_id,  miner_infos.owner_id, miner_infos.worker_id,
				miner_infos.height, miner_infos.state_root, miner_infos.storage_ask_price,
				miner_infos.min_piece_size, miner_infos.max_piece_size
				from maxht inner join miner_infos
				on miner_infos.miner_id=maxht.miner_id and miner_infos.height=maxht.maxh
				where miner_infos.miner_id > '` +
				fmt.Sprintf(*after) + `' limit ` + fmt.Sprint(*first) + `;
			`)
			if err != nil {
				fmt.Println("waserr: ", err)
			}
			defer rows.Close()
			fmt.Println("mrows", rows)
			for rows.Next() {
				err = rows.Scan(
					&minerID, &address,
					&peerID, &ownerID, &workerID,
					&height, &stateRoot,
					&storageAskPrice, &minPieceSize,
					&maxPieceSize)
				if err != nil {
					panic(err)
				}
				fmt.Println(minerID, height, peerID)
				heightInt, _ := strconv.Atoi(height)
				// minPieceSizeInt, _ := strconv.Atoi(minPieceSize)
				// maxPieceSizeInt, _ := strconv.Atoi(maxPieceSize)
				mis = append(mis, miner.MinerInfo{
					Height:  int64(heightInt),
					MinerID: minerID,
					// StateRoot:       stateRoot,
					// OwnerID:         ownerID,
					// WorkerID:        workerID,
					PeerID: peerID,
					// StorageAskPrice: storageAskPrice,
					// MinPieceSize:    uint64(minPieceSizeInt),
					// MaxPieceSize:    uint64(maxPieceSizeInt),
					Address: address,
				})
			}
		} else {
			rows, err := r.PQDB.Query(`
				with maxht as (
					select miner_id, max(height) maxh
					from miner_infos
					group by miner_id
				)
				select miner_infos.miner_id, miner_infos.address,
				miner_infos.peer_id,  miner_infos.owner_id, miner_infos.worker_id,
				miner_infos.height, miner_infos.state_root, miner_infos.storage_ask_price,
				miner_infos.min_piece_size, miner_infos.max_piece_size
				from maxht inner join miner_infos
				on miner_infos.miner_id=maxht.miner_id and miner_infos.height=maxht.maxh
				where miner_infos.miner_id > '` + fmt.Sprint(*after) + `';
			`) //, fmt.Sprintf(*after))
			if err != nil {
				fmt.Println("waserr: ", err)
			}
			defer rows.Close()
			fmt.Println("mrows", rows)
			for rows.Next() {
				err = rows.Scan(
					&minerID, &address,
					&peerID, &ownerID, &workerID,
					&height, &stateRoot,
					&storageAskPrice, &minPieceSize,
					&maxPieceSize)
				if err != nil {
					panic(err)
				}
				fmt.Println(minerID, height, peerID)
				heightInt, _ := strconv.Atoi(height)
				// minPieceSizeInt, _ := strconv.Atoi(minPieceSize)
				// maxPieceSizeInt, _ := strconv.Atoi(maxPieceSize)
				mis = append(mis, miner.MinerInfo{
					Height:  int64(heightInt),
					MinerID: minerID,
					// StateRoot:       stateRoot,
					// OwnerID:         ownerID,
					// WorkerID:        workerID,
					PeerID: peerID,
					// StorageAskPrice: storageAskPrice,
					// MinPieceSize:    uint64(minPieceSizeInt),
					// MaxPieceSize:    uint64(maxPieceSizeInt),
					Address: address,
				})
			}
		}
	} else {
		rows, err := r.PQDB.Query(`
			with maxht as (
				select miner_id, max(height) maxh
				from miner_infos
				group by miner_id
			)
			select miner_infos.miner_id, miner_infos.address,
			miner_infos.peer_id,  miner_infos.owner_id, miner_infos.worker_id,
			miner_infos.height, miner_infos.state_root, miner_infos.storage_ask_price,
			miner_infos.min_piece_size, miner_infos.max_piece_size
			from maxht inner join miner_infos
			on miner_infos.miner_id=maxht.miner_id and miner_infos.height=maxht.maxh;
		`)
		if err != nil {
			fmt.Println("waserr: ", err)
		}
		defer rows.Close()
		fmt.Println("mrows", rows)
		for rows.Next() {
			err = rows.Scan(
				&minerID, &address,
				&peerID, &ownerID, &workerID,
				&height, &stateRoot,
				&storageAskPrice, &minPieceSize,
				&maxPieceSize)
			if err != nil {
				panic(err)
			}
			fmt.Println(minerID, height, peerID)
			heightInt, _ := strconv.Atoi(height)
			// minPieceSizeInt, _ := strconv.Atoi(minPieceSize)
			// maxPieceSizeInt, _ := strconv.Atoi(maxPieceSize)
			mis = append(mis, miner.MinerInfo{
				Height:  int64(heightInt),
				MinerID: minerID,
				// StateRoot:       stateRoot,
				// OwnerID:         ownerID,
				// WorkerID:        workerID,
				PeerID: peerID,
				// StorageAskPrice: storageAskPrice,
				// MinPieceSize:    uint64(minPieceSizeInt),
				// MaxPieceSize:    uint64(maxPieceSizeInt),
				Address: address,
			})
		}
		// err := r.DB.Model(&mis).Select()
		// if err != nil {
		// 	panic(err)
		// }
	}
	fmt.Println("MIS", mis, &mis)

	var miners []*model.Miner // this is the webapp model (graphql)
	// TODO: do this in sql instead of looping again
	for _, mi := range mis {
		miners = append(miners, &model.Miner{
			ID:       mi.MinerID,
			Address:  mi.Address,
			PeerID:   mi.PeerID,
			Name:     "",
			Bio:      "",
			Location: "",
			Verified: false,
		})
	}
	return miners, nil
}

func (r *queryResolver) StorageDeal(ctx context.Context, id string) (*model.StorageDeal, error) {
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
	}
	return storagedeal, nil
}

func (r *queryResolver) AllStorageDeals(ctx context.Context, after *string, first *int, before *string, last *int, since *int, till *int) ([]*model.StorageDeal, error) {
	var mdps []market.MarketDealProposal
	if since != nil {
		if till != nil {
			err := r.DB.Model(&mdps).Where("start_epoch >= ? AND start_epoch <= ?", *since, *till).Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&mdps).Where("start_epoch >= ?", *since).Select()
			if err != nil {
				panic(err)
			}
		}
	} else {
		if till != nil {
			err := r.DB.Model(&mdps).Where("start_epoch <= ?", *till).Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&mdps).Select()
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
		})
	}
	return storagedeals, nil
}

func (r *queryResolver) Transaction(ctx context.Context, id string) (*model.Transaction, error) {
	txn := new(messages.Transaction)
	err := r.DB.Model(txn).Where("cid = ?", id).Select()
	if err != nil {
		panic(err)
	}

	transaction := &model.Transaction{
		ID:              txn.Cid,
		Amount:          txn.Amount,
		Sender:          txn.Sender,
		Receiver:        txn.Receiver,
		Height:          txn.Height,
		NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
		TransactionType: GetTransactionType(txn.MethodName),
	}
	return transaction, nil
}

func (r *queryResolver) AllTransactions(ctx context.Context, since *int, till *int) ([]*model.Transaction, error) {
	var txns []*messages.Transaction
	if since != nil {
		if till != nil {
			err := r.DB.Model(&txns).Where("height >= ? AND height <= ?", *since, *till).Select()
			if err != nil {
				panic(err)
			}
		} else {
			err := r.DB.Model(&txns).Where("height >= ?", *since).Select()
			if err != nil {
				panic(err)
			}
		}
	} else {
		if till != nil {
			err := r.DB.Model(&txns).Where("height <= ?", *till).Select()
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

	var transactions []*model.Transaction
	for _, txn := range txns {
		transactions = append(transactions, &model.Transaction{
			ID:              txn.Cid,
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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
