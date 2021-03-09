package graph

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
	"github.com/buidl-labs/filecoin-chain-indexer/model/miner"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
	"github.com/go-pg/pg/v10/orm"
)

func ComputeIncomeExpenditure(r *minerResolver, obj *model.Miner, mi *miner.MinerInfo, since *int, till *int) (*big.Int, *big.Int) {
	totalIncome := big.NewInt(0)
	totalExpenditure := big.NewInt(0)
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
	// var txns []messages.Transaction
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

				err := r.DB.Model(&ownerIncoming0).
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
						err := r.DB.Model(&ownerIncomingNotLast).
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
						err := r.DB.Model(&ownerIncomingLast).
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

				err = r.DB.Model(&workerIncoming0).
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
						err := r.DB.Model(&workerIncomingNotLast).
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, currMAC.WorkerChanges[i+1], "fil/3/storagemarket", 3, wc.From).
							Select()
						if err != nil {
							panic(err)
						}
						workerIncoming = append(workerIncoming, workerIncomingNotLast...)
					} else {
						var workerIncomingLast []messages.Transaction
						err := r.DB.Model(&workerIncomingLast).
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

				err := r.DB.Model(&ownerIncoming0).
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
						err := r.DB.Model(&ownerIncomingNotLast).
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
						err := r.DB.Model(&ownerIncomingLast).
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
				err = r.DB.Model(&workerIncoming).
					Where("height >= ? AND height <= ?", *since, *till).
					Where("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, mi.WorkerID).
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

				err := r.DB.Model(&workerIncoming0).
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
						err := r.DB.Model(&workerIncomingNotLast).
							Where("height >= ? AND height <= ?", *since, *till).
							Where("height >= ? AND height < ? AND actor_name = ? AND method = ? AND sender = ?", wc.Epoch, currMAC.WorkerChanges[i+1], "fil/3/storagemarket", 3, wc.From).
							Select()
						if err != nil {
							panic(err)
						}
						workerIncoming = append(workerIncoming, workerIncomingNotLast...)
					} else {
						var workerIncomingLast []messages.Transaction
						err := r.DB.Model(&workerIncomingLast).
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

				err = r.DB.Model(&ownerIncoming).
					Where("height >= ? AND height <= ?", *since, *till).
					WhereGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.
							WhereOr("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, mi.OwnerID).
							WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 16, obj.ID)
						return q, nil
					}).Select()
				if err != nil {
					panic(err)
				}
				incoming = append(incoming, ownerIncoming...)
			} else {
				fmt.Println("thisismycase")
				var ownerIncoming []messages.Transaction

				err := r.DB.Model(&ownerIncoming).
					Where("height >= ? AND height <= ?", *since, *till).
					WhereGroup(func(q *orm.Query) (*orm.Query, error) {
						q = q.
							WhereOr("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, mi.OwnerID).
							WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 16, obj.ID)
						return q, nil
					}).Select()
				if err != nil {
					panic(err)
				}
				fmt.Println("OIC", ownerIncoming)
				incoming = append(incoming, ownerIncoming...)

				var workerIncoming []messages.Transaction
				err = r.DB.Model(&workerIncoming).
					Where("height >= ? AND height <= ?", *since, *till).
					Where("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagemarket", 3, mi.WorkerID).
					Select()
				if err != nil {
					panic(err)
				}
				incoming = append(incoming, workerIncoming...)
			}
			// fmt.Println("ALLINCOMING", incoming)
			// end {Find incoming txns}

			// Find outgoing txns
			var outgoing0 []messages.Transaction
			err := r.DB.Model(&outgoing0).
				Where("height >= ? AND height <= ?", *since, *till).
				WhereGroup(func(q *orm.Query) (*orm.Query, error) {
					q = q.
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 5, obj.ID).
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 3, obj.ID).
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 23, obj.ID).
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 6, obj.ID).
						WhereOr("actor_name = ? AND method = ? AND receiver = ?", "fil/3/storageminer", 7, obj.ID).
						WhereOr("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagepower", 2, mi.WorkerID).
						WhereOr("actor_name = ? AND method = ? AND sender = ?", "fil/3/storagepower", 2, mi.OwnerID)
					return q, nil
				}).
				Select()
			if err != nil {
				panic(err)
			}
			outgoing = append(outgoing, outgoing0...)

			/*
				var workerOutgoing []messages.Transaction

				if workerDidChange {
					sort.Slice(currMAC.WorkerChanges[:], func(i, j int) bool {
						return currMAC.WorkerChanges[i].Epoch < currMAC.WorkerChanges[j].Epoch
					})
					wc := currMAC.WorkerChanges[0]
					var workerOutgoing0 []messages.Transaction

					err := r.DB.Model(&workerOutgoing0).
						Where("height >= ? AND height <= ?", *since, *till).
						WhereGroup(func(q *orm.Query) (*orm.Query, error) {
							q = q.
								WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 5).
								WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 3).
								WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 23).
								WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 6).
								WhereOr("height <= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.From, "fil/3/storageminer", 7).
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
							err := r.DB.Model(&workerOutgoingNotLast).
								Where("height >= ? AND height <= ?", *since, *till).
								WhereGroup(func(q *orm.Query) (*orm.Query, error) {
									q = q.
										WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 5).
										WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 3).
										WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 23).
										WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 6).
										WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer", 7).
										WhereOr("height >= ? AND height < ? AND sender = ? AND actor_name != ?", wc.Epoch, currMAC.WorkerChanges[i+1].Epoch, wc.To, "fil/3/storageminer")
									return q, nil
								}).Select()
							if err != nil {
								panic(err)
							}
							workerOutgoing = append(workerOutgoing, workerOutgoingNotLast...)
						} else {
							var workerOutgoingLast []messages.Transaction
							err := r.DB.Model(&workerOutgoingLast).
								Where("height >= ? AND height <= ?", *since, *till).
								WhereGroup(func(q *orm.Query) (*orm.Query, error) {
									q = q.
										WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 5).
										WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 3).
										WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 23).
										WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 6).
										WhereOr("height >= ? AND sender = ? AND actor_name = ? AND method != ?", wc.Epoch, wc.To, "fil/3/storageminer", 7).
										WhereOr("height >= ? AND sender = ? AND actor_name != ?", wc.Epoch, wc.To, "fil/3/storageminer")
									return q, nil
								}).Select()
							if err != nil {
								panic(err)
							}
							workerOutgoing = append(workerOutgoing, workerOutgoingLast...)
						}
					}
				} else {
					var workerOutgoing1 []messages.Transaction

					err := r.DB.Model(&workerOutgoing1).
						Where("height >= ? AND height <= ?", *since, *till).
						WhereGroup(func(q *orm.Query) (*orm.Query, error) {
							q = q.
								WhereOr("sender = ? AND actor_name = ? AND method != ?", mi.WorkerID, "fil/3/storageminer", 5).
								WhereOr("sender = ? AND actor_name = ? AND method != ?", mi.WorkerID, "fil/3/storageminer", 3).
								WhereOr("sender = ? AND actor_name = ? AND method != ?", mi.WorkerID, "fil/3/storageminer", 23).
								WhereOr("sender = ? AND actor_name = ? AND method != ?", mi.WorkerID, "fil/3/storageminer", 6).
								WhereOr("sender = ? AND actor_name = ? AND method != ?", mi.WorkerID, "fil/3/storageminer", 7).
								WhereOr("sender = ? AND actor_name != ?", mi.WorkerID, "fil/3/storageminer")

							return q, nil
						}).Select()
					if err != nil {
						panic(err)
					}

					workerOutgoing = append(workerOutgoing, workerOutgoing1...)
				}
				outgoing = append(outgoing, workerOutgoing...)
			*/

			// fmt.Println("ALLOUTGOING", outgoing)
			// end {Find outgoing txns}
		} else {
			// FIXME: the below logic needs to be changed
			// wrt changing miner addresses.
			// TODO
		}
	} else {
		if till != nil {
			// FIXME: the below logic needs to be changed
			// wrt changing miner addresses.
			// TODO
		} else {
			// FIXME: the below logic needs to be changed
			// wrt changing miner addresses.
			// TODO
		}
	}

	for _, txn := range incoming {
		amt := txn.Amount
		if txn.Amount == "0" {
			if txn.Transferred != "0" {
				amt = txn.Transferred
			}
		}
		gross := []string{amt, txn.MinerTip, txn.BaseFeeBurn}
		totalIncome = ComputeBigIntSum(totalIncome, gross)
	}
	for _, txn := range outgoing {
		amt := txn.Amount
		if txn.Amount == "0" {
			if txn.Transferred != "0" {
				amt = txn.Transferred
			}
		}
		gross := []string{amt, txn.MinerTip, txn.BaseFeeBurn}
		totalExpenditure = ComputeBigIntSum(totalExpenditure, gross)
	}
	fmt.Println("TITE", totalIncome, totalExpenditure)
	return totalIncome, totalExpenditure
}
