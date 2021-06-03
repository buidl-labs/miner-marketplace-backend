package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/buidl-labs/filecoin-chain-indexer/lens"
	"github.com/buidl-labs/miner-marketplace-backend/db/model"
	"github.com/go-pg/pg/v10"
)

func Indexer(DB *pg.DB, node lens.API) {
	hourlyTasks(DB, node)
	dailyTasks(DB, node)

	hourlyTicker := time.NewTicker(1 * time.Hour)
	dailyTicker := time.NewTicker(24 * time.Hour)

	for {
		select {
		case <-hourlyTicker.C:
			hourlyTasks(DB, node)
		case <-dailyTicker.C:
			dailyTasks(DB, node)
		}
	}
}

func hourlyTasks(DB *pg.DB, node lens.API) {
	var FILREP_MINERS string = "https://api.filrep.io/api/v1/miners"

	filRepMiners := FilRepMiners{}
	getJson(FILREP_MINERS, &filRepMiners)

	fmt.Println("pagination:", filRepMiners.Pagination)

	if filRepMiners.Pagination.Total > 0 {
		for _, m := range filRepMiners.Miners {
			reputationScoreString, _ := m.Score.(string)
			reputationScoreInt64, _ := strconv.ParseInt(reputationScoreString, 10, 64)
			reputationScore := int(reputationScoreInt64)
			// TODO! fix conversion from string to int/float
			storageAskPriceInt, _ := strconv.ParseInt(m.Price, 10, 64)
			storageAskPrice := float64(storageAskPriceInt * (10 ^ -18))
			verifiedAskPriceInt, _ := strconv.ParseInt(m.VerifiedPrice, 10, 64)
			verifiedAskPrice := float64(verifiedAskPriceInt * (10 ^ -18))

			count, _ := DB.Model((*model.Miner)(nil)).
				Where("id = ?", m.Address).
				Count()

			if count == 1 { // if already present
				var claimed bool
				err := DB.Model(&model.Miner{}).
					Column("claimed").
					Where("id = ?", m.Address).
					Select(&claimed)
				if err != nil {
					log.Println("checking if claimed:", m.Address, " error:", err)
					continue
				}
				if claimed {
					miner := &model.Miner{
						QualityAdjustedPower: m.QualityAdjPower,
						ReputationScore:      reputationScore,
					}
					_, err = DB.Model(miner).
						Column("quality_adjusted_power", "reputation_score").
						Where("id = ?", m.Address).
						Update()
					if err != nil {
						log.Println("updating miner:", m.Address, " error:", err)
						continue
					}
					fmt.Println(miner)
				} else {
					miner := &model.Miner{
						Region:               m.Region,
						Country:              m.IsoCode,
						QualityAdjustedPower: m.QualityAdjPower,
						StorageAskPrice:      storageAskPrice,
						VerifiedAskPrice:     verifiedAskPrice,
						ReputationScore:      reputationScore,
						TransparencyScore:    0,
					}
					_, err = DB.Model(miner).
						Column("region", "country", "quality_adjusted_power", "storage_ask_price", "verified_ask_price", "reputation_score", "transparency_score").
						Where("id = ?", m.Address).
						Update()
					if err != nil {
						log.Println("updating miner:", m.Address, " error:", err)
						continue
					}
				}
			} else { // if never indexed
				miner := &model.Miner{
					ID:                   m.Address,
					Claimed:              false,
					Region:               m.Region,
					Country:              m.IsoCode,
					QualityAdjustedPower: m.QualityAdjPower,
					StorageAskPrice:      storageAskPrice,
					VerifiedAskPrice:     verifiedAskPrice,
					ReputationScore:      reputationScore,
					TransparencyScore:    0,
				}
				_, err := DB.Model(miner).
					Insert()
				if err != nil {
					log.Println("inserting miner:", m.Address, " error:", err)
					continue
				}

				emptyString := ""
				minerPersonalInfo := &model.MinerPersonalInfo{
					ID:      m.Address,
					Name:    emptyString,
					Bio:     emptyString,
					Email:   emptyString,
					Website: emptyString,
					Twitter: emptyString,
					Slack:   emptyString,
				}
				_, err = DB.Model(minerPersonalInfo).Insert()
				if err != nil {
					log.Println("inserting minerPersonalInfo:", m.Address, " error:", err)
				}

				minerService := &model.MinerService{
					ID:                  m.Address,
					Storage:             true,
					Retrieval:           true,
					Repair:              false,
					DataTransferOnline:  true,
					DataTransferOffline: false,
				}
				_, err = DB.Model(minerService).Insert()
				if err != nil {
					log.Println("inserting minerService:", m.Address, " error:", err)
				}
			}
		}
	}
}

func dailyTasks(DB *pg.DB, node lens.API) {
	// update owner/worker/control addresses

	var FILREP_MINERS string = "https://api.filrep.io/api/v1/miners"
	var FILFOX_MINER string = "https://filfox.info/api/v1/address/"

	filRepMiners := new(FilRepMiners)
	getJson(FILREP_MINERS, filRepMiners)

	fmt.Println("pagination:", filRepMiners.Pagination)

	if filRepMiners.Pagination.Total > 0 {
		for _, m := range filRepMiners.Miners {
			// https://filfox.info/api/v1/address/f02770
			filFoxMiner := new(FilFoxMiner)
			getJson(FILFOX_MINER+m.Address, filFoxMiner)

			miner := &model.Miner{
				WorkerAddress: filFoxMiner.Miner.Worker.Address,
				OwnerAddress:  filFoxMiner.Miner.Owner.Address,
			}
			_, err := DB.Model(miner).
				Column("worker_address", "owner_address").
				Where("id = ?", m.Address).
				Update()
			if err != nil {
				log.Println("updating worker/owner addresses:", m.Address, " error:", err)
				continue
			}
		}
	}
}

type FilFoxMiner struct {
	Actor             string `json:"actor"`
	Address           string `json:"address"`
	Balance           string `json:"balance"`
	CreateHeight      int64  `json:"createHeight"`
	CreateTimestamp   int64  `json:"createTimestamp"`
	ID                string `json:"id"`
	LastSeenHeight    int64  `json:"lastSeenHeight"`
	LastSeenTimestamp int64  `json:"lastSeenTimestamp"`
	MessageCount      int64  `json:"messageCount"`
	Miner             struct {
		AvailableBalance string `json:"availableBalance"`
		BlocksMined      int64  `json:"blocksMined"`
		ControlAddresses []struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
		} `json:"controlAddresses"`
		InitialPledgeRequirement string      `json:"initialPledgeRequirement"`
		Location                 interface{} `json:"location"`
		MiningStats              struct {
			BlocksMined            int64   `json:"blocksMined"`
			DurationPercentage     int64   `json:"durationPercentage"`
			EquivalentMiners       int64   `json:"equivalentMiners"`
			LuckyValue             float64 `json:"luckyValue"`
			NetworkTotalRewards    string  `json:"networkTotalRewards"`
			QualityAdjPowerDelta   string  `json:"qualityAdjPowerDelta"`
			QualityAdjPowerGrowth  string  `json:"qualityAdjPowerGrowth"`
			RawBytePowerDelta      string  `json:"rawBytePowerDelta"`
			RawBytePowerGrowth     string  `json:"rawBytePowerGrowth"`
			RewardPerByte          float64 `json:"rewardPerByte"`
			TotalRewards           string  `json:"totalRewards"`
			WeightedBlocksMined    int64   `json:"weightedBlocksMined"`
			WindowedPoStFeePerByte float64 `json:"windowedPoStFeePerByte"`
		} `json:"miningStats"`
		MultiAddresses         []interface{} `json:"multiAddresses"`
		NetworkQualityAdjPower string        `json:"networkQualityAdjPower"`
		NetworkRawBytePower    string        `json:"networkRawBytePower"`
		Owner                  struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
		} `json:"owner"`
		PeerID              string `json:"peerId"`
		PledgeBalance       string `json:"pledgeBalance"`
		PreCommitDeposits   string `json:"preCommitDeposits"`
		QualityAdjPower     string `json:"qualityAdjPower"`
		QualityAdjPowerRank int64  `json:"qualityAdjPowerRank"`
		RawBytePower        string `json:"rawBytePower"`
		RawBytePowerRank    int64  `json:"rawBytePowerRank"`
		SectorPledgeBalance string `json:"sectorPledgeBalance"`
		SectorSize          int64  `json:"sectorSize"`
		Sectors             struct {
			Active     int64 `json:"active"`
			Faulty     int64 `json:"faulty"`
			Live       int64 `json:"live"`
			Recovering int64 `json:"recovering"`
		} `json:"sectors"`
		TotalRewards        string `json:"totalRewards"`
		VestingFunds        string `json:"vestingFunds"`
		WeightedBlocksMined int64  `json:"weightedBlocksMined"`
		Worker              struct {
			Address string `json:"address"`
			Balance string `json:"balance"`
		} `json:"worker"`
	} `json:"miner"`
	OwnedMiners []interface{} `json:"ownedMiners"`
	Robust      string        `json:"robust"`
	Tag         struct {
		Name   string `json:"name"`
		Signed bool   `json:"signed"`
	} `json:"tag"`
	Timestamp    int64         `json:"timestamp"`
	WorkerMiners []interface{} `json:"workerMiners"`
}

type FilRepMiners struct {
	Miners []struct {
		Address         string      `json:"address"`
		FreeSpace       string      `json:"freeSpace"`
		ID              int64       `json:"id"`
		IsoCode         string      `json:"isoCode"`
		MaxPieceSize    string      `json:"maxPieceSize"`
		MinPieceSize    string      `json:"minPieceSize"`
		Price           string      `json:"price"`
		QualityAdjPower string      `json:"qualityAdjPower"`
		Rank            string      `json:"rank"`
		RawPower        string      `json:"rawPower"`
		Region          string      `json:"region"`
		Score           interface{} `json:"score"`
		Scores          struct {
			CommittedSectorsProofs interface{} `json:"committedSectorsProofs"`
			StorageDeals           interface{} `json:"storageDeals"`
			Total                  interface{} `json:"total"`
			Uptime                 interface{} `json:"uptime"`
		} `json:"scores"`
		Status       bool `json:"status"`
		StorageDeals struct {
			AveragePrice    string `json:"averagePrice"`
			DataStored      string `json:"dataStored"`
			FaultTerminated int64  `json:"faultTerminated"`
			NoPenalties     int64  `json:"noPenalties"`
			Slashed         int64  `json:"slashed"`
			SuccessRate     string `json:"successRate"`
			Terminated      int64  `json:"terminated"`
			Total           int64  `json:"total"`
		} `json:"storageDeals"`
		UptimeAverage float64 `json:"uptimeAverage"`
		VerifiedPrice string  `json:"verifiedPrice"`
	} `json:"miners"`
	Pagination struct {
		Total int64 `json:"total"`
	} `json:"pagination"`
}

func getJson(url string, target interface{}) (interface{}, error) {
	r, err := http.Get(url)
	if err != nil {
		return target, err
	}
	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return target, readErr
	}

	jsonErr := json.Unmarshal(body, target)
	if jsonErr != nil {
		return target, jsonErr
	}
	return target, nil
}
