package dataloader

//go:generate go run github.com/vektah/dataloaden MinerLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Miner
//go:generate go run github.com/vektah/dataloaden MinerSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.Miner
//go:generate go run github.com/vektah/dataloaden StorageDealLoader int *github.com/buidl-labs/miner-marketplace-backend/graph/model.StorageDeal
//go:generate go run github.com/vektah/dataloaden StorageDealSliceLoader int []*github.com/buidl-labs/miner-marketplace-backend/graph/model.StorageDeal
//go:generate go run github.com/vektah/dataloaden TransactionLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Transaction
//go:generate go run github.com/vektah/dataloaden TransactionSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.Transaction

//go:generate go run github.com/vektah/dataloaden OwnerLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Owner
//go:generate go run github.com/vektah/dataloaden WorkerLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Worker
//go:generate go run github.com/vektah/dataloaden ContactLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Contact
//go:generate go run github.com/vektah/dataloaden ServiceDetailsLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.ServiceDetails
//go:generate go run github.com/vektah/dataloaden QualityIndicatorsLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.QualityIndicators
//go:generate go run github.com/vektah/dataloaden FinanceMetricsLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.FinanceMetrics
//go:generate go run github.com/vektah/dataloaden ServiceDetailsSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.ServiceDetails
//go:generate go run github.com/vektah/dataloaden QualityIndicatorsSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.QualityIndicators
//go:generate go run github.com/vektah/dataloaden FinanceMetricsSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.FinanceMetrics
//go:generate go run github.com/vektah/dataloaden StorageDealLoader int *github.com/buidl-labs/miner-marketplace-backend/graph/model.StorageDeal
//go:generate go run github.com/vektah/dataloaden TransactionLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Transaction
//go:generate go run github.com/vektah/dataloaden SectorLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Sector
//go:generate go run github.com/vektah/dataloaden PenaltyLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Penalty
//go:generate go run github.com/vektah/dataloaden DeadlineLoader string *github.com/buidl-labs/miner-marketplace-backend/graph/model.Deadline
//go:generate go run github.com/vektah/dataloaden StorageDealSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.StorageDeal
//go:generate go run github.com/vektah/dataloaden TransactionSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.Transaction
//go:generate go run github.com/vektah/dataloaden SectorSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.Sector
//go:generate go run github.com/vektah/dataloaden PenaltySliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.Penalty
//go:generate go run github.com/vektah/dataloaden DeadlineSliceLoader string []*github.com/buidl-labs/miner-marketplace-backend/graph/model.Deadline

// "context"
// "database/sql"
// "net/http"
// "strings"
// "time"

// "github.com/buidl-labs/miner-marketplace-backend/graph/model"

// const loadersKey = "dataloaders"

// type Loaders struct {
// 	MinerById MinerLoader
// }

// func Middleware(conn *sql.DB, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
// 			MinerById: MinerLoader{
// 				maxBatch: 100,
// 				wait:     1 * time.Millisecond,
// 				fetch: func(ids []string) ([]*model.Miner, []error) {
// 					placeholders := make([]string, len(ids))
// 					args := make([]interface{}, len(ids))
// 					for i := 0; i < len(ids); i++ {
// 						placeholders[i] = "?"
// 						args[i] = i
// 					}

// 					res := db.LogAndQuery(conn,
// 						"SELECT id, name from dataloader_example.user WHERE id IN ("+strings.Join(placeholders, ",")+")",
// 						args...,
// 					)
// 					defer res.Close()

// 					minerById := map[string]*model.Miner{}
// 					for res.Next() {
// 						miner := model.Miner{}
// 						err := res.Scan(&miner.ID, &miner.Address)
// 						if err != nil {
// 							panic(err)
// 						}
// 						minerById[miner.ID] = &miner
// 					}

// 					miners := make([]*model.Miner, len(ids))
// 					for i, id := range ids {
// 						miners[i] = minerById[id]
// 					}

// 					return miners, nil
// 				},
// 			},
// 		})
// 		r = r.WithContext(ctx)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func For(ctx context.Context) *Loaders {
// 	return ctx.Value(loadersKey).(*Loaders)
// }
