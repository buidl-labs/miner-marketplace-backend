package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/buidl-labs/filecoin-chain-indexer/model/indexing"
	"github.com/buidl-labs/filecoin-chain-indexer/model/market"
	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
	"github.com/buidl-labs/filecoin-chain-indexer/model/miner"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
)

func (r *queryResolver) ParsedTill(ctx context.Context) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Miner(ctx context.Context, id string) (*model.Miner, error) {
	pt := new(indexing.ParsedTill)
	err := r.DB.Model(pt).Limit(1).Select()
	if err != nil {
		panic(err)
	}
	fmt.Println("pth", pt.Height)
	// select * from miner_infos where id=id; (get miner from db)

	// mi := &miner.MinerInfo{MinerID: id}
	mi := new(miner.MinerInfo)
	err = r.DB.Model(mi).Where("miner_id = ? AND height = ?", id, pt.Height).Select()
	if err != nil {
		panic(err)
	}

	m := &model.Miner{
		ID:       mi.MinerID,
		Address:  mi.Address,
		PeerID:   mi.PeerID,
		Name:     "",
		Bio:      "",
		Verified: false,
	}
	// res, err := r.db.Query("")
	return m, nil
}

func (r *queryResolver) AllMiners(ctx context.Context) ([]*model.Miner, error) {
	var mis []miner.MinerInfo // this is the indexer model (db)
	err := r.DB.Model(&mis).Select()
	if err != nil {
		panic(err)
	}
	var miners []*model.Miner // this is the webapp model (graphql)
	// TODO: do this in sql instead of looping again
	for _, mi := range mis {
		miners = append(miners, &model.Miner{
			ID:       mi.MinerID,
			Address:  mi.Address,
			PeerID:   mi.PeerID,
			Name:     "",
			Bio:      "",
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

	pt := new(indexing.ParsedTill)
	err = r.DB.Model(pt).Limit(1).Select()
	if err != nil {
		panic(err)
	}
	mi := new(miner.MinerInfo)
	err = r.DB.Model(mi).Where("miner_id = ? AND height = ?", mdp.ProviderID, pt.Height).Select()
	if err != nil {
		panic(err)
	}

	m := &model.Miner{
		ID:       mi.MinerID,
		Address:  mi.Address,
		PeerID:   mi.PeerID,
		Name:     "",
		Bio:      "",
		Verified: false,
	}

	storagedeal := &model.StorageDeal{
		ID:                int(mdp.DealID),
		ClientID:          mdp.ClientID,
		StartEpoch:        mdp.StartEpoch,
		EndEpoch:          mdp.EndEpoch,
		PaddedPieceSize:   mdp.PaddedPieceSize,
		UnPaddedPieceSize: mdp.UnpaddedPieceSize,
		PieceCid:          mdp.PieceCID,
		Verified:          mdp.IsVerified,
		Miner:             m,
	}
	return storagedeal, nil
}

func (r *queryResolver) AllStorageDeals(ctx context.Context, since *int, till *int) ([]*model.StorageDeal, error) {
	var mdps []market.MarketDealProposal
	err := r.DB.Model(&mdps).Select()
	if err != nil {
		panic(err)
	}

	var storagedeals []*model.StorageDeal
	for _, mdp := range mdps {
		pt := new(indexing.ParsedTill)
		err = r.DB.Model(pt).Limit(1).Select()
		if err != nil {
			panic(err)
		}
		mi := new(miner.MinerInfo)
		err = r.DB.Model(mi).Where("miner_id = ? AND height = ?", mdp.ProviderID, pt.Height).Select()
		if err != nil {
			panic(err)
		}

		m := &model.Miner{
			ID:       mi.MinerID,
			Address:  mi.Address,
			PeerID:   mi.PeerID,
			Name:     "",
			Bio:      "",
			Verified: false,
		}

		storagedeals = append(storagedeals, &model.StorageDeal{
			ID:                int(mdp.DealID),
			ClientID:          mdp.ClientID,
			StartEpoch:        mdp.StartEpoch,
			EndEpoch:          mdp.EndEpoch,
			PaddedPieceSize:   mdp.PaddedPieceSize,
			UnPaddedPieceSize: mdp.UnpaddedPieceSize,
			PieceCid:          mdp.PieceCID,
			Verified:          mdp.IsVerified,
			Miner:             m,
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
		ID: txn.Cid,
		// Miner:           obj,
		Amount:     txn.Amount,
		Sender:     txn.Sender,
		Receiver:   txn.Receiver,
		Height:     txn.Height,
		NetworkFee: strconv.Itoa(int(txn.GasUsed)),
		// Timestamp:       time.Now(),
		TransactionType: "",
	}

	return transaction, nil
}

func (r *queryResolver) AllTransactions(ctx context.Context, since *int, till *int) ([]*model.Transaction, error) {
	var txns []*messages.Transaction
	err := r.DB.Model(&txns).Select()
	if err != nil {
		panic(err)
	}
	var transactions []*model.Transaction
	for _, txn := range txns {
		fmt.Println("some txn", txn)
		transactions = append(transactions, &model.Transaction{
			ID:              txn.Cid,
			Amount:          txn.Amount,
			Sender:          txn.Sender,
			Receiver:        txn.Receiver,
			Height:          txn.Height,
			NetworkFee:      strconv.Itoa(int(txn.GasUsed)),
			TransactionType: "",
		})
	}
	return transactions, nil
	// panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
