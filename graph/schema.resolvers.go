package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"reflect"

	dbmodel "github.com/buidl-labs/miner-marketplace-backend/db/model"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/graph/model"
	"github.com/buidl-labs/miner-marketplace-backend/service"
	"github.com/buidl-labs/miner-marketplace-backend/util"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	filecoinbig "github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/v4/actors/builtin"
	mineractor "github.com/filecoin-project/specs-actors/v4/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/v4/actors/util/smoothing"
)

func (r *aggregateEarningsResolver) Income(ctx context.Context, obj *model.AggregateEarnings) (*model.AggregateIncome, error) {
	return obj.Income, nil
}

func (r *aggregateEarningsResolver) Expenditure(ctx context.Context, obj *model.AggregateEarnings) (*model.AggregateExpenditure, error) {
	return obj.Expenditure, nil
}

func (r *aggregateEarningsResolver) NetEarnings(ctx context.Context, obj *model.AggregateEarnings) (string, error) {
	return obj.NetEarnings, nil
}

func (r *aggregateExpenditureResolver) Total(ctx context.Context, obj *model.AggregateExpenditure) (string, error) {
	return obj.Total, nil
}

func (r *aggregateExpenditureResolver) CollateralDeposit(ctx context.Context, obj *model.AggregateExpenditure) (string, error) {
	return obj.CollateralDeposit, nil
}

func (r *aggregateExpenditureResolver) Gas(ctx context.Context, obj *model.AggregateExpenditure) (string, error) {
	return obj.Gas, nil
}

func (r *aggregateExpenditureResolver) Penalty(ctx context.Context, obj *model.AggregateExpenditure) (string, error) {
	return obj.Penalty, nil
}

func (r *aggregateExpenditureResolver) Others(ctx context.Context, obj *model.AggregateExpenditure) (string, error) {
	return obj.Gas, nil
}

func (r *aggregateIncomeResolver) Total(ctx context.Context, obj *model.AggregateIncome) (string, error) {
	return obj.Total, nil
}

func (r *aggregateIncomeResolver) StorageDealPayments(ctx context.Context, obj *model.AggregateIncome) (string, error) {
	return obj.StorageDealPayments, nil
}

func (r *aggregateIncomeResolver) BlockRewards(ctx context.Context, obj *model.AggregateIncome) (string, error) {
	return obj.BlockRewards, nil
}

func (r *blockRewardsResolver) BlockRewards(ctx context.Context, obj *model.BlockRewards) (string, error) {
	return obj.BlockRewards, nil
}

func (r *blockRewardsResolver) DaysUntilEligible(ctx context.Context, obj *model.BlockRewards) (int, error) {
	return obj.DaysUntilEligible, nil
}

func (r *estimatedEarningsResolver) Income(ctx context.Context, obj *model.EstimatedEarnings) (*model.EstimatedIncome, error) {
	return obj.Income, nil
}

func (r *estimatedEarningsResolver) Expenditure(ctx context.Context, obj *model.EstimatedEarnings) (*model.EstimatedExpenditure, error) {
	return obj.Expenditure, nil
}

func (r *estimatedEarningsResolver) NetEarnings(ctx context.Context, obj *model.EstimatedEarnings) (string, error) {
	return obj.NetEarnings, nil
}

func (r *estimatedExpenditureResolver) Total(ctx context.Context, obj *model.EstimatedExpenditure) (string, error) {
	return obj.Total, nil
}

func (r *estimatedExpenditureResolver) CollateralDeposit(ctx context.Context, obj *model.EstimatedExpenditure) (string, error) {
	return obj.CollateralDeposit, nil
}

func (r *estimatedExpenditureResolver) Gas(ctx context.Context, obj *model.EstimatedExpenditure) (string, error) {
	return obj.Gas, nil
}

func (r *estimatedExpenditureResolver) Penalty(ctx context.Context, obj *model.EstimatedExpenditure) (string, error) {
	return obj.Penalty, nil
}

func (r *estimatedExpenditureResolver) Others(ctx context.Context, obj *model.EstimatedExpenditure) (string, error) {
	return obj.Others, nil
}

func (r *estimatedIncomeResolver) Total(ctx context.Context, obj *model.EstimatedIncome) (string, error) {
	return obj.Total, nil
}

func (r *estimatedIncomeResolver) StorageDealPayments(ctx context.Context, obj *model.EstimatedIncome) (*model.StorageDealPayments, error) {
	return obj.StorageDealPayments, nil
}

func (r *estimatedIncomeResolver) BlockRewards(ctx context.Context, obj *model.EstimatedIncome) (*model.BlockRewards, error) {
	return obj.BlockRewards, nil
}

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

func (r *minerResolver) StorageDealStats(ctx context.Context, obj *model.Miner) (*model.StorageDealStats, error) {
	minerStorageDealStats := dbmodel.MinerStorageDealStats{}
	err := r.DB.Model(&minerStorageDealStats).Where("id = ?", obj.ID).Select()
	if err != nil {
		return &model.StorageDealStats{}, nil
	}
	return &model.StorageDealStats{
		AveragePrice:    minerStorageDealStats.AveragePrice,
		DataStored:      minerStorageDealStats.DataStored,
		FaultTerminated: int(minerStorageDealStats.FaultTerminated),
		NoPenalties:     int(minerStorageDealStats.NoPenalties),
		Slashed:         int(minerStorageDealStats.Slashed),
		SuccessRate:     minerStorageDealStats.SuccessRate,
		Terminated:      int(minerStorageDealStats.Terminated),
		Total:           int(minerStorageDealStats.Total),
	}, nil
}

func (r *minerResolver) Transactions(ctx context.Context, obj *model.Miner) ([]*model.Transaction, error) {
	var dbTransactions []*dbmodel.Transaction
	if err := r.DB.Model(&dbTransactions).
		Where("miner_id = ?", obj.ID).
		Select(); err != nil {
		return []*model.Transaction{}, err
	}
	var transactions []*model.Transaction
	for _, dbTransaction := range dbTransactions {
		transactions = append(transactions, &model.Transaction{
			ID:              dbTransaction.ID,
			Miner:           obj,
			Height:          int(dbTransaction.Height),
			TransactionType: dbTransaction.TransactionType,
			MethodName:      dbTransaction.MethodName,
			Value:           dbTransaction.Value,
			MinerFee:        dbTransaction.MinerFee,
			BurnFee:         dbTransaction.BurnFee,
			From:            dbTransaction.From,
			To:              dbTransaction.To,
			ExitCode:        dbTransaction.ExitCode,
			Deals:           dbTransaction.Deals,
		})
	}
	return transactions, nil
}

func (r *minerResolver) AggregateEarnings(ctx context.Context, obj *model.Miner, startHeight int, endHeight int, transactionTypes []bool, includeGas bool) (*model.AggregateEarnings, error) {
	useAllTransactionTypes := false
	if transactionTypes == nil {
		fmt.Println("nil", transactionTypes)
		useAllTransactionTypes = true
	} else if len(transactionTypes) == 0 {
		fmt.Println("empty", transactionTypes)
		useAllTransactionTypes = true
	} else if len(transactionTypes) != 6 {
		return &model.AggregateEarnings{
			Income:      &model.AggregateIncome{Total: "0", StorageDealPayments: "0", BlockRewards: "0"},
			Expenditure: &model.AggregateExpenditure{Total: "0", CollateralDeposit: "0", Gas: "0", Penalty: "0", Others: "0"},
			NetEarnings: "0",
		}, fmt.Errorf("length of transactionTypes array should be 6")
	}

	fmt.Println("useAllTransactionTypes", useAllTransactionTypes)

	var dbTransactions []*dbmodel.Transaction
	if useAllTransactionTypes {
		if err := r.DB.Model(&dbTransactions).
			Where("miner_id = ?", obj.ID).
			Where("height >= ?", startHeight).
			Where("height <= ?", endHeight).
			Where("exit_code = ?", 0).
			Select(); err != nil {
			return &model.AggregateEarnings{
				Income:      &model.AggregateIncome{Total: "0", StorageDealPayments: "0", BlockRewards: "0"},
				Expenditure: &model.AggregateExpenditure{Total: "0", CollateralDeposit: "0", Gas: "0", Penalty: "0", Others: "0"},
				NetEarnings: "0",
			}, err
		}
	} else {
		transactionTypesQuery := util.GenerateTransactionTypesQuery(transactionTypes)
		if err := r.DB.Model(&dbTransactions).
			Where("miner_id = ?", obj.ID).
			Where("height >= ?", startHeight).
			Where("height <= ?", endHeight).
			Where(transactionTypesQuery).
			Where("exit_code = ?", 0).
			Select(); err != nil {
			return &model.AggregateEarnings{
				Income:      &model.AggregateIncome{Total: "0", StorageDealPayments: "0", BlockRewards: "0"},
				Expenditure: &model.AggregateExpenditure{Total: "0", CollateralDeposit: "0", Gas: "0", Penalty: "0", Others: "0"},
				NetEarnings: "0",
			}, err
		}
	}

	storageDealPayments := big.NewInt(0)

	var dbMarketDealProposals []*dbmodel.MarketDealProposal
	err := r.DB.Model(&dbMarketDealProposals).
		Where("provider = ?", obj.ID).
		Where("start_epoch <= ?", endHeight).
		Where("end_epoch >= ?", startHeight).
		Select()

	if err == nil {
		for _, mdp := range dbMarketDealProposals {
			fmt.Println("###### dealid:", mdp.ID)
			duration := mdp.EndEpoch - mdp.StartEpoch
			durationBigInt := big.NewInt(duration)
			storagePrice := new(big.Int)
			storagePrice, ok := storagePrice.SetString(mdp.StoragePrice, 10)
			if !ok {
				fmt.Println("SetString: error")
			}
			fmt.Println("storagePrice", storagePrice)
			pricePerEpoch := new(big.Int).Div(storagePrice, durationBigInt)
			fmt.Println("pricePerEpoch", pricePerEpoch)
			fmt.Println("durationBigInt", durationBigInt)
			var startEpochInSelectedRange int64
			var endEpochInSelectedRange int64
			if mdp.StartEpoch > int64(startHeight) {
				startEpochInSelectedRange = mdp.StartEpoch
			} else {
				startEpochInSelectedRange = int64(startHeight)
			}
			if mdp.EndEpoch < int64(endHeight) {
				endEpochInSelectedRange = mdp.EndEpoch
			} else {
				endEpochInSelectedRange = int64(endHeight)
			}
			earningsFromCurrentDeal := new(big.Int).Mul(pricePerEpoch, big.NewInt(endEpochInSelectedRange-startEpochInSelectedRange))
			fmt.Println("earningsFromCurrentDeal", earningsFromCurrentDeal)
			storageDealPayments = new(big.Int).Add(storageDealPayments, earningsFromCurrentDeal)
			fmt.Println("################")
		}
	}

	income := big.NewInt(0)
	blockRewards := big.NewInt(0)

	expenditure := big.NewInt(0)
	collateralDeposit := big.NewInt(0)
	gas := big.NewInt(0)
	penalty := big.NewInt(0)
	others := big.NewInt(0)

	for _, dbTransaction := range dbTransactions {
		switch dbTransaction.MethodName {
		case "PreCommitSector", "ProveCommitSector":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			collateralDeposit = new(big.Int).Add(collateralDeposit, val)
			expenditure = new(big.Int).Add(expenditure, val)
			minerFee, ok := new(big.Int).SetString(dbTransaction.MinerFee, 10)
			if !ok {
				fmt.Println("problem converting minerFee to bigint:", dbTransaction.MinerFee, "id:", dbTransaction.ID)
			}
			burnFee, ok := new(big.Int).SetString(dbTransaction.BurnFee, 10)
			if !ok {
				fmt.Println("problem converting burnFee to bigint:", dbTransaction.BurnFee, "id:", dbTransaction.ID)
			}
			gas = new(big.Int).Add(gas, minerFee)
			gas = new(big.Int).Add(gas, burnFee)
			if includeGas {
				collateralDeposit = new(big.Int).Add(collateralDeposit, minerFee)
				collateralDeposit = new(big.Int).Add(collateralDeposit, burnFee)
				expenditure = new(big.Int).Add(expenditure, minerFee)
				expenditure = new(big.Int).Add(expenditure, burnFee)
			}
		case "TerminateSectors", "RepayDebt",
			"ReportConsensusFault", "DisputeWindowedPoSt":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			penalty = new(big.Int).Add(penalty, val)
			expenditure = new(big.Int).Add(expenditure, val)
			minerFee, ok := new(big.Int).SetString(dbTransaction.MinerFee, 10)
			if !ok {
				fmt.Println("problem converting minerFee to bigint:", dbTransaction.MinerFee, "id:", dbTransaction.ID)
			}
			burnFee, ok := new(big.Int).SetString(dbTransaction.BurnFee, 10)
			if !ok {
				fmt.Println("problem converting burnFee to bigint:", dbTransaction.BurnFee, "id:", dbTransaction.ID)
			}
			gas = new(big.Int).Add(gas, minerFee)
			gas = new(big.Int).Add(gas, burnFee)
			if includeGas {
				penalty = new(big.Int).Add(penalty, minerFee)
				penalty = new(big.Int).Add(penalty, burnFee)
				expenditure = new(big.Int).Add(expenditure, minerFee)
				expenditure = new(big.Int).Add(expenditure, burnFee)
			}
		case "SubmitWindowedPoSt", "ChangeWorkerAddress", "ChangePeerID",
			"ExtendSectorExpiration", "DeclareFaults", "DeclareFaultsRecovered",
			"ChangeMultiaddrs", "CompactSectorNumbers", "ConfirmUpdateWorkerKey",
			"ChangeOwnerAddress":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			others = new(big.Int).Add(others, val)
			expenditure = new(big.Int).Add(expenditure, val)
			minerFee, ok := new(big.Int).SetString(dbTransaction.MinerFee, 10)
			if !ok {
				fmt.Println("problem converting minerFee to bigint:", dbTransaction.MinerFee, "id:", dbTransaction.ID)
			}
			burnFee, ok := new(big.Int).SetString(dbTransaction.BurnFee, 10)
			if !ok {
				fmt.Println("problem converting burnFee to bigint:", dbTransaction.BurnFee, "id:", dbTransaction.ID)
			}
			gas = new(big.Int).Add(gas, minerFee)
			gas = new(big.Int).Add(gas, burnFee)
			if includeGas {
				others = new(big.Int).Add(others, minerFee)
				others = new(big.Int).Add(others, burnFee)
				expenditure = new(big.Int).Add(expenditure, minerFee)
				expenditure = new(big.Int).Add(expenditure, burnFee)
			}
		case "ApplyRewards":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			blockRewards = new(big.Int).Add(blockRewards, val)
			income = new(big.Int).Add(income, val)
		}
	}

	income = new(big.Int).Add(income, storageDealPayments)
	netEarnings := new(big.Int).Sub(income, expenditure)

	fmt.Println("income", income, "expenditure", expenditure, "netEarnings", netEarnings)

	return &model.AggregateEarnings{
		Income: &model.AggregateIncome{
			Total:               income.String(),
			StorageDealPayments: storageDealPayments.String(),
			BlockRewards:        blockRewards.String(),
		},
		Expenditure: &model.AggregateExpenditure{
			Total:             expenditure.String(),
			CollateralDeposit: collateralDeposit.String(),
			Gas:               gas.String(),
			Penalty:           penalty.String(),
			Others:            others.String(),
		},
		NetEarnings: netEarnings.String(),
	}, nil
}

func (r *minerResolver) EstimatedEarnings(ctx context.Context, obj *model.Miner, days int, transactionTypes []bool, includeGas bool) (*model.EstimatedEarnings, error) {
	existingStorageDealPayments := big.NewInt(0)
	potentialFutureDealPayments := big.NewInt(0)

	ts, _ := r.LensAPI.ChainHead(context.Background())

	startHeight := ts.Height()
	endHeight := startHeight + 60*2880
	var dbMarketDealProposals []*dbmodel.MarketDealProposal
	err := r.DB.Model(&dbMarketDealProposals).
		Where("provider = ?", obj.ID).
		Where("start_epoch <= ?", endHeight).
		Where("end_epoch >= ?", startHeight).
		Select()

	if err == nil {
		for _, mdp := range dbMarketDealProposals {
			fmt.Println("###### dealid:", mdp.ID)
			duration := mdp.EndEpoch - mdp.StartEpoch
			durationBigInt := big.NewInt(duration)
			storagePrice := new(big.Int)
			storagePrice, ok := storagePrice.SetString(mdp.StoragePrice, 10)
			if !ok {
				fmt.Println("SetString: error")
			}
			fmt.Println("storagePrice", storagePrice)
			pricePerEpoch := new(big.Int).Div(storagePrice, durationBigInt)
			fmt.Println("pricePerEpoch", pricePerEpoch)
			fmt.Println("durationBigInt", durationBigInt)
			var startEpochInSelectedRange int64
			var endEpochInSelectedRange int64
			if mdp.StartEpoch > int64(startHeight) {
				startEpochInSelectedRange = mdp.StartEpoch
			} else {
				startEpochInSelectedRange = int64(startHeight)
			}
			if mdp.EndEpoch < int64(endHeight) {
				endEpochInSelectedRange = mdp.EndEpoch
			} else {
				endEpochInSelectedRange = int64(endHeight)
			}
			earningsFromCurrentDeal := new(big.Int).Mul(pricePerEpoch, big.NewInt(endEpochInSelectedRange-startEpochInSelectedRange))
			fmt.Println("earningsFromCurrentDeal", earningsFromCurrentDeal)
			existingStorageDealPayments = new(big.Int).Add(existingStorageDealPayments, earningsFromCurrentDeal)
			fmt.Println("################")
		}
	}

	pricePerEpochSum := big.NewInt(0)
	var last2MonthsMarketDealProposals []*dbmodel.MarketDealProposal
	err = r.DB.Model(&last2MonthsMarketDealProposals).
		Where("provider = ?", obj.ID).
		Where("start_epoch >= ?", ts.Height()-60*2880).
		Select()
	if err == nil {
		// estimatedFutureDealsEarnings := int64(0)
		last2MonthsDealsCount := int64(len(last2MonthsMarketDealProposals))
		if last2MonthsDealsCount != 0 {
			for _, mdp := range last2MonthsMarketDealProposals {
				duration := mdp.EndEpoch - mdp.StartEpoch
				durationBigInt := big.NewInt(duration)
				storagePrice := new(big.Int)
				storagePrice, ok := storagePrice.SetString(mdp.StoragePrice, 10)
				if !ok {
					fmt.Println("SetString: error")
				}
				fmt.Println("storagePrice", storagePrice)
				pricePerEpoch := new(big.Int).Div(storagePrice, durationBigInt)

				// storagePricePerEpoch, _ := strconv.ParseInt(mdp.StoragePricePerEpoch, 10, 64)
				pricePerEpochSum = new(big.Int).Add(pricePerEpochSum, pricePerEpoch)
			}
			averagePricePerEpochLast2Months := new(big.Int).Div(pricePerEpochSum, big.NewInt(last2MonthsDealsCount))
			// estimatedFutureDealsFloat := new(big.Float).Mul(
			// 	new(big.Float).Quo(
			// 		big.NewFloat(float64(last2MonthsDealsCount)),
			// 		big.NewFloat(float64(60))),
			// 	big.NewFloat(float64(days)),
			// )
			// estimatedFutureDeals := estimatedFutureDealsFloat.Int
			estimatedFutureDeals := new(big.Int).Mul(
				new(big.Int).Div(
					big.NewInt(last2MonthsDealsCount),
					big.NewInt(int64(60))),
				big.NewInt(int64(days)),
			)
			potentialFutureDealPayments = new(big.Int).Mul(
				new(big.Int).Mul(
					estimatedFutureDeals,
					averagePricePerEpochLast2Months,
				),
				big.NewInt(int64(days*2880)),
			)
			// averagePricePerEpochLast2Months := pricePerEpochSum / last2MonthsDealsCount
			// estimatedFutureDeals := (last2MonthsDealsCount / 60) * int64(days)
			// estimatedFutureDealsEarnings = estimatedFutureDeals * averagePricePerEpochLast2Months * int64(days*2880)
		}
	}

	dbTransaction := dbmodel.Transaction{}
	err = r.DB.Model(&dbTransaction).
		Where("method_name = 'CreateMiner'").
		Where("miner_id = ?", obj.ID).
		Select()
	if err != nil {
		return &model.EstimatedEarnings{
			Income: &model.EstimatedIncome{
				Total: new(big.Int).Add(existingStorageDealPayments, potentialFutureDealPayments).String(),
				StorageDealPayments: &model.StorageDealPayments{
					ExistingDeals:        existingStorageDealPayments.String(),
					PotentialFutureDeals: potentialFutureDealPayments.String(),
				}, BlockRewards: &model.BlockRewards{
					BlockRewards:      "0",
					DaysUntilEligible: 0,
				}},
			Expenditure: &model.EstimatedExpenditure{Total: "0", CollateralDeposit: "0", Gas: "0", Penalty: "0", Others: "0"},
			NetEarnings: new(big.Int).Add(existingStorageDealPayments, potentialFutureDealPayments).String(),
		}, err
	}
	minerCreationHeight := int(dbTransaction.Height)

	var dbTransactions []*dbmodel.Transaction
	// transactionTypesQuery := util.GenerateTransactionTypesQuery(transactionTypes)
	if err := r.DB.Model(&dbTransactions).
		Where("miner_id = ?", obj.ID).
		// Where("height >= ?", startHeight).
		// Where("height <= ?", endHeight).
		// Where(transactionTypesQuery).
		Where("exit_code = ?", 0).
		Select(); err != nil {
		return &model.EstimatedEarnings{
			Income: &model.EstimatedIncome{
				Total: new(big.Int).Add(existingStorageDealPayments, potentialFutureDealPayments).String(),
				StorageDealPayments: &model.StorageDealPayments{
					ExistingDeals:        existingStorageDealPayments.String(),
					PotentialFutureDeals: potentialFutureDealPayments.String(),
				}, BlockRewards: &model.BlockRewards{
					BlockRewards:      "0",
					DaysUntilEligible: 0,
				}},
			Expenditure: &model.EstimatedExpenditure{Total: "0", CollateralDeposit: "0", Gas: "0", Penalty: "0", Others: "0"},
			NetEarnings: new(big.Int).Add(existingStorageDealPayments, potentialFutureDealPayments).String(),
		}, err
	}

	income := big.NewInt(0)
	// storageDealPayments := big.NewInt(0) // TODO: estimate storageDealPayments
	blockRewards := big.NewInt(0)

	expenditure := big.NewInt(0)
	collateralDeposit := big.NewInt(0)
	gas := big.NewInt(0)
	penalty := big.NewInt(0)
	others := big.NewInt(0)

	for _, dbTransaction := range dbTransactions {
		switch dbTransaction.MethodName {
		case "PreCommitSector", "ProveCommitSector":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			collateralDeposit = new(big.Int).Add(collateralDeposit, val)
			expenditure = new(big.Int).Add(expenditure, val)
			minerFee, ok := new(big.Int).SetString(dbTransaction.MinerFee, 10)
			if !ok {
				fmt.Println("problem converting minerFee to bigint:", dbTransaction.MinerFee, "id:", dbTransaction.ID)
			}
			burnFee, ok := new(big.Int).SetString(dbTransaction.BurnFee, 10)
			if !ok {
				fmt.Println("problem converting burnFee to bigint:", dbTransaction.BurnFee, "id:", dbTransaction.ID)
			}
			gas = new(big.Int).Add(gas, minerFee)
			gas = new(big.Int).Add(gas, burnFee)
			if includeGas {
				collateralDeposit = new(big.Int).Add(collateralDeposit, minerFee)
				collateralDeposit = new(big.Int).Add(collateralDeposit, burnFee)
				expenditure = new(big.Int).Add(expenditure, minerFee)
				expenditure = new(big.Int).Add(expenditure, burnFee)
			}
		case "TerminateSectors", "RepayDebt",
			"ReportConsensusFault", "DisputeWindowedPoSt":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			penalty = new(big.Int).Add(penalty, val)
			expenditure = new(big.Int).Add(expenditure, val)
			minerFee, ok := new(big.Int).SetString(dbTransaction.MinerFee, 10)
			if !ok {
				fmt.Println("problem converting minerFee to bigint:", dbTransaction.MinerFee, "id:", dbTransaction.ID)
			}
			burnFee, ok := new(big.Int).SetString(dbTransaction.BurnFee, 10)
			if !ok {
				fmt.Println("problem converting burnFee to bigint:", dbTransaction.BurnFee, "id:", dbTransaction.ID)
			}
			gas = new(big.Int).Add(gas, minerFee)
			gas = new(big.Int).Add(gas, burnFee)
			if includeGas {
				penalty = new(big.Int).Add(penalty, minerFee)
				penalty = new(big.Int).Add(penalty, burnFee)
				expenditure = new(big.Int).Add(expenditure, minerFee)
				expenditure = new(big.Int).Add(expenditure, burnFee)
			}
		case "SubmitWindowedPoSt", "ChangeWorkerAddress", "ChangePeerID",
			"ExtendSectorExpiration", "DeclareFaults", "DeclareFaultsRecovered",
			"ChangeMultiaddrs", "CompactSectorNumbers", "ConfirmUpdateWorkerKey",
			"ChangeOwnerAddress":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			others = new(big.Int).Add(others, val)
			expenditure = new(big.Int).Add(expenditure, val)
			minerFee, ok := new(big.Int).SetString(dbTransaction.MinerFee, 10)
			if !ok {
				fmt.Println("problem converting minerFee to bigint:", dbTransaction.MinerFee, "id:", dbTransaction.ID)
			}
			burnFee, ok := new(big.Int).SetString(dbTransaction.BurnFee, 10)
			if !ok {
				fmt.Println("problem converting burnFee to bigint:", dbTransaction.BurnFee, "id:", dbTransaction.ID)
			}
			gas = new(big.Int).Add(gas, minerFee)
			gas = new(big.Int).Add(gas, burnFee)
			if includeGas {
				others = new(big.Int).Add(others, minerFee)
				others = new(big.Int).Add(others, burnFee)
				expenditure = new(big.Int).Add(expenditure, minerFee)
				expenditure = new(big.Int).Add(expenditure, burnFee)
			}
		case "ApplyRewards":
			val, ok := new(big.Int).SetString(dbTransaction.Value, 10)
			if !ok {
				fmt.Println("problem converting value to bigint:", dbTransaction.Value, "id:", dbTransaction.ID)
			}
			blockRewards = new(big.Int).Add(blockRewards, val)
			// income = new(big.Int).Add(income, val)
		}
	}

	expenditurePerDay := new(big.Int).Div(expenditure, big.NewInt(int64(minerCreationHeight)))
	collateralDepositPerDay := new(big.Int).Div(collateralDeposit, big.NewInt(int64(minerCreationHeight)))
	gasPerDay := new(big.Int).Div(gas, big.NewInt(int64(minerCreationHeight)))
	penaltyPerDay := new(big.Int).Div(penalty, big.NewInt(int64(minerCreationHeight)))
	othersPerDay := new(big.Int).Div(others, big.NewInt(int64(minerCreationHeight)))

	estimatedExpenditure := new(big.Int).Mul(expenditurePerDay, big.NewInt(int64(days)))
	estimatedCollateralDeposit := new(big.Int).Mul(collateralDepositPerDay, big.NewInt(int64(days)))
	estimatedGas := new(big.Int).Mul(gasPerDay, big.NewInt(int64(days)))
	estimatedPenalty := new(big.Int).Mul(penaltyPerDay, big.NewInt(int64(days)))
	estimatedOthers := new(big.Int).Mul(othersPerDay, big.NewInt(int64(days)))

	minerID, _ := address.NewFromString(obj.ID)
	powerActorID, _ := address.NewFromString("f04")
	rewardActorID, _ := address.NewFromString("f02")
	ts, _ = r.LensAPI.ChainHead(context.Background())

	daysUntilEligible := big.NewInt(-1)
	nrwd := filecoinbig.NewInt(0)
	minQAP := big.NewInt(10995116277760) // 10 TiB
	minerPower, _ := r.LensAPI.StateMinerPower(context.Background(), minerID, ts.Key())
	cmpR := minerPower.MinerPower.QualityAdjPower.Int.Cmp(minQAP)
	if cmpR == -1 {
		// find daysUntilEligible
		lastMonthTs, _ := r.LensAPI.ChainGetTipSetByHeight(context.Background(), ts.Height()-30*2880, types.EmptyTSK)
		minerPowerLastMonth, _ := r.LensAPI.StateMinerPower(context.Background(), minerID, lastMonthTs.Key())
		dailyPowerGrowthLastMonth := new(big.Int).Div(
			new(big.Int).Sub(
				minerPower.MinerPower.QualityAdjPower.Int,
				minerPowerLastMonth.MinerPower.QualityAdjPower.Int,
			),
			big.NewInt(int64(30)),
		)
		fmt.Println("dailyPowerGrowthLastMonth", dailyPowerGrowthLastMonth)
		if dailyPowerGrowthLastMonth.String() == "0" {
			fmt.Println("minerPower.MinerPower.QualityAdjPower.Int", minerPower.MinerPower.QualityAdjPower.Int)
			if minerPower.MinerPower.QualityAdjPower.Int.Cmp(big.NewInt(0)) == 1 { // if miner's current QAP=0
				daysUntilEligible = new(big.Int).Div(
					new(big.Int).Sub(
						minerPower.MinerPower.QualityAdjPower.Int,
						big.NewInt(0),
					),
					big.NewInt(int64(minerCreationHeight)), // TODO: replace 30 with days since miner was created
				)
			}
		} else {
			daysUntilEligible = new(big.Int).Div(
				new(big.Int).Sub(minQAP, minerPower.MinerPower.QualityAdjPower.Int),
				dailyPowerGrowthLastMonth,
			)
		}
	} else {
		daysUntilEligible = big.NewInt(0)
		PowerActorState, err := r.LensAPI.StateReadState(context.Background(), powerActorID, ts.Key())
		if err != nil {
			panic(err)
		}
		RewardActorState, err := r.LensAPI.StateReadState(context.Background(), rewardActorID, ts.Key())
		if err != nil {
			panic(err)
		}

		pas, _ := PowerActorState.State.(map[string]interface{})
		ras, _ := RewardActorState.State.(map[string]interface{})
		// fmt.Println(reflect.TypeOf(pas["ThisEpochQAPowerSmoothed"]), " ", pas["ThisEpochQAPowerSmoothed"])

		ThisEpochQAPowerSmoothed, _ := pas["ThisEpochQAPowerSmoothed"].(map[string]interface{})
		// fmt.Println(reflect.TypeOf(ThisEpochQAPowerSmoothed), ThisEpochQAPowerSmoothed,
		// 	"pe:", ThisEpochQAPowerSmoothed["PositionEstimate"],
		// 	"ve:", ThisEpochQAPowerSmoothed["VelocityEstimate"])

		ThisEpochRewardSmoothed, _ := ras["ThisEpochRewardSmoothed"].(map[string]interface{})
		// fmt.Println(reflect.TypeOf(ThisEpochRewardSmoothed), ThisEpochRewardSmoothed,
		// 	"pe:", ThisEpochRewardSmoothed["PositionEstimate"],
		// 	"ve:", ThisEpochRewardSmoothed["VelocityEstimate"])

		a := ThisEpochQAPowerSmoothed["PositionEstimate"].(string)
		ThisEpochQAPowerSmoothedPositionEstimate, _ := new(big.Int).SetString(a, 10)

		b := ThisEpochQAPowerSmoothed["VelocityEstimate"].(string)
		ThisEpochQAPowerSmoothedVelocityEstimate, _ := new(big.Int).SetString(b, 10)

		c := ThisEpochRewardSmoothed["PositionEstimate"].(string)
		ThisEpochRewardSmoothedPositionEstimate, _ := new(big.Int).SetString(c, 10)

		d := ThisEpochRewardSmoothed["VelocityEstimate"].(string)
		ThisEpochRewardSmoothedVelocityEstimate, _ := new(big.Int).SetString(d, 10)

		// fmt.Println("pas", pas, " old ", reflect.TypeOf(PowerActorState.State), " ", PowerActorState.State)
		// fmt.Println("ras", ras, " old ", reflect.TypeOf(RewardActorState.State), " ", RewardActorState.State)

		qaPower := minerPower.MinerPower.QualityAdjPower // filecoinbig.NewInt(int64(100000 * math.Pow(2, 30)))
		fmt.Println("minerqaPower", qaPower)
		nrwd = mineractor.ExpectedRewardForPower(smoothing.FilterEstimate{
			PositionEstimate: filecoinbig.NewFromGo(ThisEpochRewardSmoothedPositionEstimate),
			VelocityEstimate: filecoinbig.NewFromGo(ThisEpochRewardSmoothedVelocityEstimate),
		}, smoothing.FilterEstimate{
			PositionEstimate: filecoinbig.NewFromGo(ThisEpochQAPowerSmoothedPositionEstimate),
			VelocityEstimate: filecoinbig.NewFromGo(ThisEpochQAPowerSmoothedVelocityEstimate),
		}, qaPower, builtin.EpochsInDay*abi.ChainEpoch(days))

		// atto := big.NewInt(1e18)
		// minerProjectedReward := nrwd.Int.Div(nrwd.Int, atto)
		// fmt.Println("minerProjectedReward", minerProjectedReward)
	}

	fmt.Println("nrwd", nrwd)
	fmt.Println("BEFORE daysUntilEligible", daysUntilEligible)
	daysUntilEligibleInt := 0
	if daysUntilEligible.Cmp(big.NewInt(-1)) == 0 && cmpR == -1 {
		daysUntilEligibleInt = int(math.Inf(1))
		fmt.Println("miner power zero, daysUntilEligibleInt", daysUntilEligibleInt)
	} else {
		daysUntilEligibleInt = int(daysUntilEligible.Int64())
		fmt.Println("else, daysUntilEligibleInt", daysUntilEligibleInt)
	}
	fmt.Println("NOW daysUntilEligible", daysUntilEligible)

	income = new(big.Int).Add(existingStorageDealPayments, potentialFutureDealPayments)
	income = new(big.Int).Add(income, nrwd.Int)
	netEarnings := new(big.Int).Sub(income, estimatedExpenditure)

	return &model.EstimatedEarnings{
		Income: &model.EstimatedIncome{
			Total: income.String(),
			StorageDealPayments: &model.StorageDealPayments{
				ExistingDeals:        existingStorageDealPayments.String(),
				PotentialFutureDeals: potentialFutureDealPayments.String(),
			},
			BlockRewards: &model.BlockRewards{
				BlockRewards:      nrwd.String(),
				DaysUntilEligible: daysUntilEligibleInt,
			},
		},
		Expenditure: &model.EstimatedExpenditure{
			Total:             estimatedExpenditure.String(),
			CollateralDeposit: estimatedCollateralDeposit.String(),
			Gas:               estimatedGas.String(),
			Penalty:           estimatedPenalty.String(),
			Others:            estimatedOthers.String(),
		},
		NetEarnings: netEarnings.String(),
	}, nil
}

func (r *mutationResolver) ClaimProfile(ctx context.Context, input model.ProfileClaimInput) (bool, error) {
	fmt.Println("i", input.MinerID, "t", reflect.TypeOf(input.MinerID))
	fmt.Println("j", input.LedgerAddress, "t", reflect.TypeOf(input.LedgerAddress))

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

	filFoxStatsPower := new(service.FilFoxStatsPower)
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

func (r *serviceResolver) ServiceTypes(ctx context.Context, obj *model.Service) (*model.ServiceTypes, error) {
	return obj.ServiceTypes, nil
}

func (r *serviceResolver) DataTransferMechanism(ctx context.Context, obj *model.Service) (*model.DataTransferMechanism, error) {
	return obj.DataTransferMechanism, nil
}

func (r *storageDealPaymentsResolver) ExistingDeals(ctx context.Context, obj *model.StorageDealPayments) (string, error) {
	return obj.ExistingDeals, nil
}

func (r *storageDealPaymentsResolver) PotentialFutureDeals(ctx context.Context, obj *model.StorageDealPayments) (string, error) {
	return obj.PotentialFutureDeals, nil
}

func (r *transactionResolver) Miner(ctx context.Context, obj *model.Transaction) (*model.Miner, error) {
	return obj.Miner, nil
}

func (r *workerResolver) Miner(ctx context.Context, obj *model.Worker) (*model.Miner, error) {
	return obj.Miner, nil
}

// AggregateEarnings returns generated.AggregateEarningsResolver implementation.
func (r *Resolver) AggregateEarnings() generated.AggregateEarningsResolver {
	return &aggregateEarningsResolver{r}
}

// AggregateExpenditure returns generated.AggregateExpenditureResolver implementation.
func (r *Resolver) AggregateExpenditure() generated.AggregateExpenditureResolver {
	return &aggregateExpenditureResolver{r}
}

// AggregateIncome returns generated.AggregateIncomeResolver implementation.
func (r *Resolver) AggregateIncome() generated.AggregateIncomeResolver {
	return &aggregateIncomeResolver{r}
}

// BlockRewards returns generated.BlockRewardsResolver implementation.
func (r *Resolver) BlockRewards() generated.BlockRewardsResolver { return &blockRewardsResolver{r} }

// EstimatedEarnings returns generated.EstimatedEarningsResolver implementation.
func (r *Resolver) EstimatedEarnings() generated.EstimatedEarningsResolver {
	return &estimatedEarningsResolver{r}
}

// EstimatedExpenditure returns generated.EstimatedExpenditureResolver implementation.
func (r *Resolver) EstimatedExpenditure() generated.EstimatedExpenditureResolver {
	return &estimatedExpenditureResolver{r}
}

// EstimatedIncome returns generated.EstimatedIncomeResolver implementation.
func (r *Resolver) EstimatedIncome() generated.EstimatedIncomeResolver {
	return &estimatedIncomeResolver{r}
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

// StorageDealPayments returns generated.StorageDealPaymentsResolver implementation.
func (r *Resolver) StorageDealPayments() generated.StorageDealPaymentsResolver {
	return &storageDealPaymentsResolver{r}
}

// Transaction returns generated.TransactionResolver implementation.
func (r *Resolver) Transaction() generated.TransactionResolver { return &transactionResolver{r} }

// Worker returns generated.WorkerResolver implementation.
func (r *Resolver) Worker() generated.WorkerResolver { return &workerResolver{r} }

type aggregateEarningsResolver struct{ *Resolver }
type aggregateExpenditureResolver struct{ *Resolver }
type aggregateIncomeResolver struct{ *Resolver }
type blockRewardsResolver struct{ *Resolver }
type estimatedEarningsResolver struct{ *Resolver }
type estimatedExpenditureResolver struct{ *Resolver }
type estimatedIncomeResolver struct{ *Resolver }
type locationResolver struct{ *Resolver }
type minerResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type ownerResolver struct{ *Resolver }
type personalInfoResolver struct{ *Resolver }
type pricingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
type storageDealPaymentsResolver struct{ *Resolver }
type transactionResolver struct{ *Resolver }
type workerResolver struct{ *Resolver }
