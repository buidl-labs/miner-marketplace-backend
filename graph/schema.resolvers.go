package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	dbmodel "github.com/buidl-labs/miner-marketplace-backend/db/model"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
	"github.com/buidl-labs/miner-marketplace-backend/service"
	"github.com/buidl-labs/miner-marketplace-backend/util"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	// "github.com/alecthomas/units"
)

func (r *locationResolver) Region(ctx context.Context, obj *model.Location) (string, error) {
	return obj.Region, nil
}

func (r *locationResolver) Country(ctx context.Context, obj *model.Location) (string, error) {
	return obj.Country, nil
}

func (r *minerResolver) PersonalInfo(ctx context.Context, obj *model.Miner) (*model.PersonalInfo, error) {
	minerPersonalInfo := dbmodel.MinerPersonalInfo{}
	err := r.DB.Model(&minerPersonalInfo).Where("id = ?", obj.ID).Select()
	if err != nil {
		return &model.PersonalInfo{}, nil
	}
	return &model.PersonalInfo{
		Name:    minerPersonalInfo.Name,
		Bio:     minerPersonalInfo.Bio,
		Email:   minerPersonalInfo.Email,
		Website: minerPersonalInfo.Website,
		Twitter: minerPersonalInfo.Twitter,
		Slack:   minerPersonalInfo.Slack,
	}, nil
}

func (r *minerResolver) Worker(ctx context.Context, obj *model.Miner) (*model.Worker, error) {
	dbMiner := dbmodel.Miner{}
	if err := r.DB.Model(&dbMiner).Where("id = ?", obj.ID).Select(); err != nil {
		return &model.Worker{}, err
	}
	return &model.Worker{
		ID:      dbMiner.WorkerID,
		Address: dbMiner.WorkerAddress,
	}, nil
}

func (r *minerResolver) Owner(ctx context.Context, obj *model.Miner) (*model.Owner, error) {
	dbMiner := dbmodel.Miner{}
	if err := r.DB.Model(&dbMiner).Where("id = ?", obj.ID).Select(); err != nil {
		return &model.Owner{}, err
	}
	return &model.Owner{
		ID:      dbMiner.OwnerID,
		Address: dbMiner.OwnerAddress,
	}, nil
}

func (r *minerResolver) Location(ctx context.Context, obj *model.Miner) (*model.Location, error) {
	dbMiner := dbmodel.Miner{}
	if err := r.DB.Model(&dbMiner).Where("id = ?", obj.ID).Select(); err != nil {
		return &model.Location{}, err
	}
	return &model.Location{
		Region:  dbMiner.Region,
		Country: dbMiner.Country,
	}, nil
}

func (r *minerResolver) QualityAdjustedPower(ctx context.Context, obj *model.Miner) (string, error) {
	return obj.QualityAdjustedPower, nil
}

func (r *minerResolver) Service(ctx context.Context, obj *model.Miner) (*model.Service, error) {
	dbMinerService := dbmodel.MinerService{}
	if err := r.DB.Model(&dbMinerService).Where("id = ?", obj.ID).Select(); err != nil {
		return &model.Service{}, err
	}
	return &model.Service{
		ServiceTypes: &model.ServiceTypes{
			Storage:   dbMinerService.Storage,
			Retrieval: dbMinerService.Retrieval,
			Repair:    dbMinerService.Repair,
		},
		DataTransferMechanism: &model.DataTransferMechanism{
			Online:  dbMinerService.DataTransferOnline,
			Offline: dbMinerService.DataTransferOffline,
		},
	}, nil
}

func (r *minerResolver) Pricing(ctx context.Context, obj *model.Miner) (*model.Pricing, error) {
	dbMiner := dbmodel.Miner{}
	if err := r.DB.Model(&dbMiner).Where("id = ?", obj.ID).Select(); err != nil {
		fmt.Println("Pricing: ", err)
		return &model.Pricing{}, err
	}
	return &model.Pricing{
		StorageAskPrice:   dbMiner.StorageAskPrice,
		VerifiedAskPrice:  dbMiner.VerifiedAskPrice,
		RetrievalAskPrice: dbMiner.RetrievalAskPrice,
	}, nil
}

func (r *minerResolver) ReputationScore(ctx context.Context, obj *model.Miner) (int, error) {
	return obj.ReputationScore, nil
}

func (r *minerResolver) TransparencyScore(ctx context.Context, obj *model.Miner) (int, error) {
	return obj.TransparencyScore, nil
}

func (r *mutationResolver) ClaimProfile(ctx context.Context, input model.ProfileClaimInput) (bool, error) {
	fmt.Println("i", input.MinerID, "t", reflect.TypeOf(input.MinerID))
	fmt.Println("j", input.LedgerAddress, "t", reflect.TypeOf(input.LedgerAddress))

	var isClaimed bool
	err := r.DB.Model((*dbmodel.Miner)(nil)).
		Column("claimed").
		Where("id = ?", input.MinerID).
		Select(&isClaimed)

	// ######
	// NOTE: just for testing with our ledger wallets
	if input.MinerID == "f04321" {
		if input.LedgerAddress == "f1v2qntmt4k6wxugdbxqjw6l3knywh2csi2lcmz2a" ||
			input.LedgerAddress == "f1rb4xvch25rqshc7oklj3wcxgotezciqbjufgeli" ||
			input.LedgerAddress == "f1zi7hgjoxpbfci3s5ggiexnwoi2c6gsnu74agt7a" {
			dbMiner := dbmodel.Miner{
				Claimed:           true,
				TransparencyScore: 10,
			}
			_, err := r.DB.Model(&dbMiner).
				Column("claimed", "transparency_score").
				Where("id = ?", input.MinerID).
				Update()
			if err != nil {
				fmt.Println("testing.", err)
				return false, err // failed to update in db
			}
			return true, nil
		} else {
			return false, nil
		}
	}
	// ######

	minerID, err := address.NewFromString(input.MinerID)
	if err != nil {
		fmt.Println("1.", err)
		return false, err
	}
	minerInfo, _ := r.LensAPI.StateMinerInfo(context.Background(), minerID, types.EmptyTSK)
	if err != nil {
		fmt.Println("2.", err)
		return false, err
	}
	ownerAddress, err := r.LensAPI.StateAccountKey(context.Background(), minerInfo.Owner, types.EmptyTSK)
	if err != nil {
		fmt.Println("3.", err)
		return false, err
	}

	fmt.Println("cmp", "la", input.LedgerAddress, "oa", ownerAddress.String())

	if input.LedgerAddress == ownerAddress.String() {
		// success
		dbMiner := dbmodel.Miner{
			Claimed:           true,
			TransparencyScore: 10,
		}
		_, err := r.DB.Model(&dbMiner).
			Column("claimed", "transparency_score").
			Where("id = ?", input.MinerID).
			Update()
		if err != nil {
			fmt.Println("4.", err)
			return false, err // failed to update in db
		}
		return true, nil
	} else {
		// failure
		return false, nil
	}
}

func (r *mutationResolver) EditProfile(ctx context.Context, input model.ProfileSettingsInput) (bool, error) {
	dbMiner := dbmodel.Miner{}
	if err := r.DB.Model(&dbMiner).Where("id = ?", input.MinerID).Select(); err != nil {
		return false, err
	}
	if dbMiner.Claimed {
		// ######
		// NOTE: just for testing with our ledger wallets
		if input.MinerID == "f04321" {
			if input.LedgerAddress == "f1v2qntmt4k6wxugdbxqjw6l3knywh2csi2lcmz2a" ||
				input.LedgerAddress == "f1rb4xvch25rqshc7oklj3wcxgotezciqbjufgeli" ||
				input.LedgerAddress == "f1zi7hgjoxpbfci3s5ggiexnwoi2c6gsnu74agt7a" {
				updatedMiner := dbmodel.Miner{
					Region:            input.Region,
					Country:           input.Country,
					StorageAskPrice:   input.StorageAskPrice,
					VerifiedAskPrice:  input.VerifiedAskPrice,
					RetrievalAskPrice: input.RetrievalAskPrice,
					TransparencyScore: service.ComputeTransparencyScore(input),
				}
				updatedMinerPersonalInfo := dbmodel.MinerPersonalInfo{
					Name:    input.Name,
					Bio:     input.Bio,
					Email:   input.Email,
					Website: input.Website,
					Twitter: input.Twitter,
					Slack:   input.Slack,
				}
				updatedMinerService := dbmodel.MinerService{
					Storage:             input.Storage,
					Retrieval:           input.Retrieval,
					Repair:              input.Repair,
					DataTransferOnline:  input.Online,
					DataTransferOffline: input.Offline,
				}

				_, err := r.DB.Model(&updatedMiner).
					Column("region", "country", "storage_ask_price", "verified_ask_price", "retrieval_ask_price", "transparency_score").
					Where("id = ?", input.MinerID).
					Update()
				if err != nil {
					return false, err
				}

				_, err = r.DB.Model(&updatedMinerPersonalInfo).
					Column("name", "bio", "email", "website", "twitter", "slack").
					Where("id = ?", input.MinerID).
					Update()
				if err != nil {
					return false, err
				}

				_, err = r.DB.Model(&updatedMinerService).
					Column("storage", "retrieval", "repair", "data_transfer_online", "data_transfer_offline").
					Where("id = ?", input.MinerID).
					Update()
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
		// ######
		minerID, err := address.NewFromString(input.MinerID)
		if err != nil {
			return false, err
		}
		minerInfo, _ := r.LensAPI.StateMinerInfo(context.Background(), minerID, types.EmptyTSK)
		if err != nil {
			return false, err
		}
		ownerAddress, err := r.LensAPI.StateAccountKey(context.Background(), minerInfo.Owner, types.EmptyTSK)
		if err != nil {
			return false, err
		}
		if input.LedgerAddress == ownerAddress.String() {
			updatedMiner := dbmodel.Miner{
				Region:            input.Region,
				Country:           input.Country,
				StorageAskPrice:   input.StorageAskPrice,
				VerifiedAskPrice:  input.VerifiedAskPrice,
				RetrievalAskPrice: input.RetrievalAskPrice,
				TransparencyScore: service.ComputeTransparencyScore(input),
			}
			updatedMinerPersonalInfo := dbmodel.MinerPersonalInfo{
				Name:    input.Name,
				Bio:     input.Bio,
				Email:   input.Email,
				Website: input.Website,
				Twitter: input.Twitter,
				Slack:   input.Slack,
			}
			updatedMinerService := dbmodel.MinerService{
				Storage:             input.Storage,
				Retrieval:           input.Retrieval,
				Repair:              input.Repair,
				DataTransferOnline:  input.Online,
				DataTransferOffline: input.Offline,
			}

			_, err := r.DB.Model(&updatedMiner).
				Column("region", "country", "storage_ask_price", "verified_ask_price", "retrieval_ask_price", "transparency_score").
				Where("id = ?", input.MinerID).
				Update()
			if err != nil {
				return false, err
			}

			_, err = r.DB.Model(&updatedMinerPersonalInfo).
				Column("name", "bio", "email", "website", "twitter", "slack").
				Where("id = ?", input.MinerID).
				Update()
			if err != nil {
				return false, err
			}

			_, err = r.DB.Model(&updatedMinerService).
				Column("storage", "retrieval", "repair", "data_transfer_online", "data_transfer_offline").
				Where("id = ?", input.MinerID).
				Update()
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (r *ownerResolver) Miner(ctx context.Context, obj *model.Owner) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *personalInfoResolver) Name(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	return obj.Name, nil
}

func (r *personalInfoResolver) Bio(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	return obj.Bio, nil
}

func (r *personalInfoResolver) Email(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	return obj.Email, nil
}

func (r *personalInfoResolver) Website(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	return obj.Website, nil
}

func (r *personalInfoResolver) Twitter(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	return obj.Twitter, nil
}

func (r *personalInfoResolver) Slack(ctx context.Context, obj *model.PersonalInfo) (string, error) {
	return obj.Slack, nil
}

func (r *pricingResolver) StorageAskPrice(ctx context.Context, obj *model.Pricing) (string, error) {
	return obj.StorageAskPrice, nil
}

func (r *pricingResolver) VerifiedAskPrice(ctx context.Context, obj *model.Pricing) (string, error) {
	return obj.VerifiedAskPrice, nil
}

func (r *pricingResolver) RetrievalAskPrice(ctx context.Context, obj *model.Pricing) (string, error) {
	return obj.RetrievalAskPrice, nil
}

func (r *queryResolver) Miner(ctx context.Context, id string) (*model.Miner, error) {
	dbMiner := dbmodel.Miner{}
	if err := r.DB.Model(&dbMiner).Where("id = ?", id).Select(); err != nil {
		fmt.Println("Miner", err)
		return &model.Miner{}, err
	}
	return &model.Miner{
		ID:                   dbMiner.ID,
		Claimed:              dbMiner.Claimed,
		QualityAdjustedPower: dbMiner.QualityAdjustedPower,
		ReputationScore:      dbMiner.ReputationScore,
		TransparencyScore:    dbMiner.TransparencyScore,
	}, nil
}

func (r *queryResolver) Miners(ctx context.Context, first *int, offset *int) ([]*model.Miner, error) {
	var _first = 100
	var _offset = 0

	if first != nil {
		_first = *first
	}
	if offset != nil {
		_offset = *offset
	}

	var dbMiners []*dbmodel.Miner
	if err := r.DB.Model(&dbMiners).
		Limit(_first).
		Offset(_offset).
		Select(); err != nil {
		return []*model.Miner{}, err
	}
	var miners []*model.Miner
	for _, dbMiner := range dbMiners {
		miners = append(miners, &model.Miner{
			ID:                   dbMiner.ID,
			Claimed:              dbMiner.Claimed,
			QualityAdjustedPower: dbMiner.QualityAdjustedPower,
			ReputationScore:      dbMiner.ReputationScore,
			TransparencyScore:    dbMiner.TransparencyScore,
		})
	}
	return miners, nil
}

func (r *queryResolver) NetworkStats(ctx context.Context) (*model.NetworkStats, error) {
	var dbMiners []*dbmodel.Miner
	var activeMinersCount int
	activeMinersCount, err := r.DB.Model(&dbMiners).Count()
	if err != nil {
		fmt.Println("couldn't get activeminerscount")
		activeMinersCount = 2247
	}
	var FILFOX_STATS_POWER string = "https://filfox.info/api/v1/stats/power"

	filFoxStatsPower := new(FilFoxStatsPower)
	util.GetJson(FILFOX_STATS_POWER, filFoxStatsPower)

	// fmt.Println("pagination:")
	// for _, fsp := range *filFoxStatsPower {
	// 	fmt.Println(fsp.QualityAdjPower)
	// }
	fsp := *filFoxStatsPower
	fmt.Println("zeroth", fsp[0].QualityAdjPower, fsp[0])
	var powerEiB string = "6.131 EiB"
	n := new(big.Int)
	n, ok := n.SetString(fsp[0].QualityAdjPower, 10)
	if !ok {
		fmt.Println("SetString: error")
	}
	fmt.Println("n", n, "factor", big.NewInt(int64(1.153e+18)))
	// big.
	dv := new(big.Int).Div(n, big.NewInt(int64(1.153e+18)))
	fmt.Println("dv", dv)
	// powerEiB = dv.String()
	fmt.Println("dvs", powerEiB)
	// 1.153e+18
	return &model.NetworkStats{
		ActiveMinersCount:      activeMinersCount,
		NetworkStorageCapacity: powerEiB,
		DataStored:             powerEiB,
	}, nil
}

type FilFoxStatsPower []struct {
	Height               int    `json:"height"`
	Timestamp            int    `json:"timestamp"`
	RawBytePower         string `json:"rawBytePower"`
	QualityAdjPower      string `json:"qualityAdjPower"`
	RawBytePowerDelta    string `json:"rawBytePowerDelta"`
	QualityAdjPowerDelta string `json:"qualityAdjPowerDelta"`
}

func (r *serviceResolver) ServiceTypes(ctx context.Context, obj *model.Service) (*model.ServiceTypes, error) {
	return obj.ServiceTypes, nil
}

func (r *serviceResolver) DataTransferMechanism(ctx context.Context, obj *model.Service) (*model.DataTransferMechanism, error) {
	return obj.DataTransferMechanism, nil
}

func (r *workerResolver) Miner(ctx context.Context, obj *model.Worker) (*model.Miner, error) {
	return obj.Miner, nil
}

// Location returns generated.LocationResolver implementation.
func (r *Resolver) Location() generated.LocationResolver { return &locationResolver{r} }

// Miner returns generated.MinerResolver implementation.
func (r *Resolver) Miner() generated.MinerResolver { return &minerResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Owner returns generated.OwnerResolver implementation.
func (r *Resolver) Owner() generated.OwnerResolver { return &ownerResolver{r} }

// PersonalInfo returns generated.PersonalInfoResolver implementation.
func (r *Resolver) PersonalInfo() generated.PersonalInfoResolver { return &personalInfoResolver{r} }

// Pricing returns generated.PricingResolver implementation.
func (r *Resolver) Pricing() generated.PricingResolver { return &pricingResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Service returns generated.ServiceResolver implementation.
func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

// Worker returns generated.WorkerResolver implementation.
func (r *Resolver) Worker() generated.WorkerResolver { return &workerResolver{r} }

type locationResolver struct{ *Resolver }
type minerResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type ownerResolver struct{ *Resolver }
type personalInfoResolver struct{ *Resolver }
type pricingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
type workerResolver struct{ *Resolver }
