package service

import (
	// "context"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/buidl-labs/filecoin-chain-indexer/lens"
	"github.com/buidl-labs/miner-marketplace-backend/db/model"
	gqlmodel "github.com/buidl-labs/miner-marketplace-backend/graph/model"
	"github.com/buidl-labs/miner-marketplace-backend/util"
	"github.com/spf13/viper"
	// "github.com/ipfs/go-cid"
	// "github.com/filecoin-project/go-address"
	// "github.com/filecoin-project/lotus/api"
	// "github.com/filecoin-project/lotus/chain/types"
	"github.com/go-pg/pg/v10"
)

var yamlExample = []byte(`
miners:
- f0115238
- f08403
`)

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
	util.GetJson(FILREP_MINERS, &filRepMiners)

	fmt.Println("pagination:", filRepMiners.Pagination)

	if filRepMiners.Pagination.Total > 0 {
		for _, m := range filRepMiners.Miners {
			reputationScoreString, _ := m.Score.(string)
			reputationScoreInt64, _ := strconv.ParseInt(reputationScoreString, 10, 64)
			reputationScore := int(reputationScoreInt64)

			storageAskPrice := m.Price
			verifiedAskPrice := m.VerifiedPrice

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
	util.GetJson(FILREP_MINERS, filRepMiners)

	fmt.Println("pagination:", filRepMiners.Pagination)

	if filRepMiners.Pagination.Total > 0 {
		for _, m := range filRepMiners.Miners {
			// https://filfox.info/api/v1/address/f02770
			filFoxMiner := new(FilFoxMiner)
			util.GetJson(FILFOX_MINER+m.Address, filFoxMiner)

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

func ComputeTransparencyScore(input gqlmodel.ProfileSettingsInput) int {
	transparencyScore := 10.0 // already claimed
	if input.Name != "" {
		transparencyScore += 5
	}
	if input.Bio != "" {
		transparencyScore += 5
	}
	if input.Slack != "" {
		transparencyScore += 15
	}
	if input.Twitter != "" {
		transparencyScore += 15
	}
	if input.Email != "" {
		transparencyScore += 7.5
	}
	if input.Website != "" {
		transparencyScore += 7.5
	}
	transparencyScore += 10 // for service details, give all points for datatransfermechanism and servicetype
	if input.Region != "" {
		transparencyScore += 2.5
	}
	if input.Country != "" {
		transparencyScore += 2.5
	}
	if input.StorageAskPrice != "" {
		transparencyScore += 20
	}
	return int(transparencyScore)
}

func MinerPageMessages(DB *pg.DB, node lens.API) {
	/*
		fmt.Println("MinerPageMessages")
		rewardActor, _ := address.NewFromString("f02")
		fmt.Println("rewardActor", rewardActor) // From: rewardActor,
		miner, _ := address.NewFromString("f0123261")
		fmt.Println("miner", miner)
		ts, err := node.ChainGetTipSetByHeight(context.Background(), 828482, types.EmptyTSK)
		if err != nil {
			fmt.Println("CGTBH err:", err)
		}
		fmt.Println("ts", ts, "tsk", ts.Key())
		cids, err := node.StateListMessages(context.Background(), &api.MessageMatch{From: rewardActor, To: miner}, types.EmptyTSK, 828482)
		if err != nil {
			fmt.Println("SLMs err:", err)
		}
		fmt.Println("cids", cids)
		for _, cid := range cids {
			msg, err := node.ChainGetMessage(context.Background(), cid)
			if err != nil {
				fmt.Println("CGM err:", err)
			}
			fmt.Println("msg", msg)
		}
	*/

	// reward
	// https://filfox.info/api/v1/address/f0152712/transfers?pageSize=100&page=0&type=reward

	var FILFOX_MINER string = "https://filfox.info/api/v1/address/"
	var FILFOX_MESSAGE string = "https://filfox.info/api/v1/message/"

	// var FILREP_MINER_TRANSFERS="https://filfox.info/api/v1/address/f0152712/transfers?pageSize=100&page=0&type=reward"

	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	miners := viper.GetStringSlice("miners")

	for _, m := range miners {
		fmt.Println("miner", m)
		// block rewards
		minerRewards := []struct {
			// type AutoGenerated struct {
			// TotalCount int `json:"totalCount"`
			// Blocks     []struct {
			Cid          string `json:"cid"`
			Height       int    `json:"height"`
			Timestamp    int    `json:"timestamp"`
			Size         int    `json:"size"`
			WinCount     int    `json:"winCount"`
			Reward       string `json:"reward"`
			Penalty      string `json:"penalty"`
			MessageCount int    `json:"messageCount"`
			// } `json:"blocks"`
			// }
			// Height    int    `json:"height"`
			// Timestamp int    `json:"timestamp"`
			// From      string `json:"from"`
			// FromTag   struct {
			// 	Name   string `json:"name"`
			// 	Signed bool   `json:"signed"`
			// } `json:"fromTag"`
			// To    string `json:"to"`
			// ToTag struct {
			// 	Name   string `json:"name"`
			// 	Signed bool   `json:"signed"`
			// } `json:"toTag"`
			// Value string `json:"value"`
			// Type  string `json:"type"`
		}{}
		filFoxMinerTransfersReward := new(FilFoxMinerBlocks) // new(FilFoxMinerTransfers)
		// https://filfox.info/api/v1/address/f062353/blocks?pageSize=100&page=0
		util.GetJson(FILFOX_MINER+m+"/blocks?pageSize=100&page=0", filFoxMinerTransfersReward)

		// fmt.Println("url:", FILFOX_MINER+m+"/transfers?pageSize=100&page=0&type=reward")
		// util.GetJson(FILFOX_MINER+m+"/transfers?pageSize=100&page=0&type=reward", filFoxMinerTransfersReward)
		// fmt.Println("latest reward page", filFoxMinerTransfersReward)
		minerRewards = append(minerRewards, filFoxMinerTransfersReward.Blocks...)

		var db_miner_transfers_reward_total_count int64
		err := DB.Model((*model.FilfoxMinerMessagesCount)(nil)).
			Column("miner_transfers_reward_total_count").
			Where("id = ?", m).
			Select(&db_miner_transfers_reward_total_count)

		fmt.Println("db_miner_transfers_reward_total_count", db_miner_transfers_reward_total_count)

		totalRewardCount := filFoxMinerTransfersReward.TotalCount
		fmt.Println("totalRewardCount", totalRewardCount)

		var diffRw int64
		var pagesRw int
		if err == nil && db_miner_transfers_reward_total_count < int64(totalRewardCount) {
			diffRw = int64(totalRewardCount) - db_miner_transfers_reward_total_count
			pagesRw = int(diffRw) / 100
			fmt.Println("case1 diffRw", diffRw, "pagesRw", pagesRw)
			for i := 1; i <= pagesRw; i++ {
				fmt.Println("page", i)
				fmt.Println("iterminerRewards", len(minerRewards))
				util.GetJson(FILFOX_MINER+m+"/blocks?pageSize=100&page="+fmt.Sprintf("%d", i), filFoxMinerTransfersReward)
				minerRewards = append(minerRewards, filFoxMinerTransfersReward.Blocks...)

				// util.GetJson(FILFOX_MINER+m+"/transfers?pageSize=100&page="+fmt.Sprintf("%d", i)+"&type=reward", filFoxMinerTransfersReward)
				// minerRewards = append(minerRewards, filFoxMinerTransfersReward.Transfers...)

				// util.GetJson(FILFOX_MINER+m+"/messages?pageSize=100&page="+fmt.Sprintf("%d", i), filFoxMinerTransfersReward)
			}
		} else if db_miner_transfers_reward_total_count != int64(totalRewardCount) {
			minerRewardPagesCount := totalRewardCount / 100
			fmt.Println("minerRewardPagesCount", minerRewardPagesCount)
			for i := 1; i <= minerRewardPagesCount; i++ {
				fmt.Println("page", i)
				fmt.Println("iterminerRewards", len(minerRewards))
				util.GetJson(FILFOX_MINER+m+"/blocks?pageSize=100&page="+fmt.Sprintf("%d", i), filFoxMinerTransfersReward)
				minerRewards = append(minerRewards, filFoxMinerTransfersReward.Blocks...)
				// util.GetJson(FILFOX_MINER+m+"/transfers?pageSize=100&page="+fmt.Sprintf("%d", i)+"&type=reward", filFoxMinerTransfersReward)
				// minerRewards = append(minerRewards, filFoxMinerTransfersReward.Transfers...)
			}
		}
		// DB.Model(&model.FilfoxMinerMessagesCount{
		// 	MinerTransfersRewardTotalCount: int64(totalRewardCount),
		// }).
		// 	Column("miner_transfers_reward_total_count").
		// 	Where("id = ?", m).
		// 	Update()
		fmt.Println("minerRewards", len(minerRewards))

		if db_miner_transfers_reward_total_count != int64(totalRewardCount) {
			for _, mr := range minerRewards {
				value := "0"
				if mr.Reward != "" {
					value = mr.Reward
				}
				_, err := DB.Model(&model.Transaction{
					ID:              mr.Cid, // mr.To + fmt.Sprintf("%v", mr.Height) + "reward", // cid not available in filfox
					MinerID:         m,
					Height:          int64(mr.Height),
					Timestamp:       int64(mr.Timestamp),
					TransactionType: "Block Reward",
					MethodName:      "ApplyRewards",
					Value:           value,
					MinerFee:        "0",
					BurnFee:         "0",
					From:            "f02",
					To:              m,
					ExitCode:        0,
				}).Insert()
				if err != nil {
					fmt.Println("minerRewards insert err:", err)
				}
			}
		}
		_, err = DB.Model(&model.FilfoxMinerMessagesCount{
			// ID:                             m,
			MinerTransfersRewardTotalCount: int64(totalRewardCount),
		}).
			// OnConflict("(id) DO UPDATE").
			// Set("title = EXCLUDED.title").
			Column("miner_transfers_reward_total_count").
			Where("id = ?", m).
			Update()
		if err != nil {
			fmt.Println("inserting/updating MinerTransfersRewardTotalCount", err)
		}

		// miner actor messages
		// https://filfox.info/api/v1/address/f0115238/messages?pageSize=100&page=0
		minerMessages := []struct {
			Cid       string `json:"cid"`
			Height    int    `json:"height"`
			Timestamp int    `json:"timestamp"`
			From      string `json:"from"`
			To        string `json:"to"`
			Nonce     int    `json:"nonce"`
			Value     string `json:"value"`
			Method    string `json:"method"`
			Receipt   struct {
				ExitCode int `json:"exitCode"`
			} `json:"receipt"`
		}{}
		filFoxMinerMessages := new(FilFoxMinerMessages)
		fmt.Println("url:", FILFOX_MINER+m+"/messages?pageSize=100&page=0")
		util.GetJson(FILFOX_MINER+m+"/messages?pageSize=100&page=0", filFoxMinerMessages)
		// fmt.Println("latest miner messages page", filFoxMinerMessages)
		minerMessages = append(minerMessages, filFoxMinerMessages.Messages...)

		var db_miner_messages_total_count int64
		err = DB.Model((*model.FilfoxMinerMessagesCount)(nil)).
			Column("miner_messages_total_count").
			Where("id = ?", m).
			Select(&db_miner_messages_total_count)

		fmt.Println("db_miner_messages_total_count", db_miner_messages_total_count)

		totalMinerMessageCount := filFoxMinerMessages.TotalCount
		fmt.Println("totalMinerMessageCount", totalMinerMessageCount)

		var diff int64
		var pages int
		if err == nil && db_miner_messages_total_count < int64(totalMinerMessageCount) {
			diff = int64(totalMinerMessageCount) - db_miner_messages_total_count
			pages = int(diff) / 100
			fmt.Println("case1 diff", diff, "pages", pages)
			for i := 1; i <= pages; i++ {
				fmt.Println("page", i)
				fmt.Println("iterminerMessages", len(minerMessages))
				util.GetJson(FILFOX_MINER+m+"/messages?pageSize=100&page="+fmt.Sprintf("%d", i), filFoxMinerMessages)
				minerMessages = append(minerMessages, filFoxMinerMessages.Messages...)
			}
		} else if db_miner_messages_total_count != int64(totalMinerMessageCount) {
			minerMessagePagesCount := totalMinerMessageCount / 100
			fmt.Println("minerMessagePagesCount", minerMessagePagesCount)
			for i := 1; i <= minerMessagePagesCount; i++ {
				fmt.Println("page", i)
				fmt.Println("iterminerMessages", len(minerMessages))
				util.GetJson(FILFOX_MINER+m+"/messages?pageSize=100&page="+fmt.Sprintf("%d", i), filFoxMinerMessages)
				minerMessages = append(minerMessages, filFoxMinerMessages.Messages...)
			}
		}
		// DB.Model(&model.FilfoxMinerMessagesCount{
		// 	MinerMessagesTotalCount: int64(totalMinerMessageCount),
		// }).
		// 	Column("miner_messages_total_count").
		// 	Where("id = ?", m).
		// 	Update()
		// fmt.Println("minerMessages", len(minerMessages))
		if db_miner_messages_total_count != int64(totalMinerMessageCount) {
			for _, mr := range minerMessages {
				// https://filfox.info/api/v1/message/bafy2bzacebo54zcaakbqov2e7shpfvxqugmpgmn4m7mpirsbac6w7jumkra3i
				filFoxMessage := new(FilFoxMessage)
				util.GetJson(FILFOX_MESSAGE+mr.Cid, filFoxMessage)
				transactionType, value, minerFee, burnFee, _, _ := GetMessageAttributes(node, *filFoxMessage)
				if value == "" {
					value = "0"
				}
				if minerFee == "" {
					minerFee = "0"
				}
				if burnFee == "" {
					burnFee = "0"
				}
				_, err := DB.Model(&model.Transaction{
					ID:              mr.Cid,
					MinerID:         m,
					Height:          int64(mr.Height),
					Timestamp:       int64(mr.Timestamp),
					TransactionType: transactionType,
					MethodName:      mr.Method,
					Value:           value,
					MinerFee:        minerFee,
					BurnFee:         burnFee,
					From:            mr.From,
					To:              mr.To,
					ExitCode:        mr.Receipt.ExitCode,
				}).Insert()
				if err != nil {
					fmt.Println("minerMessages insert err:", err)
				}
			}
		}
		_, err = DB.Model(&model.FilfoxMinerMessagesCount{
			// ID:                      m,
			MinerMessagesTotalCount: int64(totalMinerMessageCount),
		}).
			// OnConflict("(id) DO UPDATE").
			// Set("title = EXCLUDED.title").
			Column("miner_messages_total_count").
			Where("id = ?", m).
			Update()
		if err != nil {
			fmt.Println("inserting/updating MinerMessagesTotalCount", err)
		}
	}
}

func PublishStorageDealsMessages(DB *pg.DB, node lens.API) {
	var FILFOX_MESSAGE string = "https://filfox.info/api/v1/message/"

	filFoxMessagesList := new(FilFoxMessagesList)
	util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page=0&method=PublishStorageDeals", filFoxMessagesList)

	publishStorageDealsMessages := []struct {
		Cid       string `json:"cid"`
		Height    int    `json:"height"`
		Timestamp int    `json:"timestamp"`
		From      string `json:"from"`
		To        string `json:"to"`
		ToTag     struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag"`
		Value   string `json:"value"`
		Method  string `json:"method"`
		Receipt struct {
			ExitCode int `json:"exitCode"`
		} `json:"receipt"`
		FromTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"fromTag,omitempty"`
	}{}

	publishStorageDealsMessages = append(publishStorageDealsMessages, filFoxMessagesList.Messages...)

	var db_publish_storage_deals_messages_total_count int64
	err := DB.Model((*model.FilfoxMessagesCount)(nil)).
		Column("publish_storage_deals_messages_total_count").
		Where("id = 'psdm'").
		Select(&db_publish_storage_deals_messages_total_count)

	fmt.Println("db_publish_storage_deals_messages_total_count", db_publish_storage_deals_messages_total_count)

	totalPublishStorageDealsMessageCount := filFoxMessagesList.TotalCount
	fmt.Println("totalPublishStorageDealsMessageCount", totalPublishStorageDealsMessageCount)

	var diff int64
	var pages int
	if err == nil && db_publish_storage_deals_messages_total_count < int64(totalPublishStorageDealsMessageCount) {
		diff = int64(totalPublishStorageDealsMessageCount) - db_publish_storage_deals_messages_total_count
		pages = int(diff) / 100
		fmt.Println("case1 diff", diff, "pages", pages)
		for i := 1; i <= pages; i++ {
			fmt.Println("page", i)
			fmt.Println("iterminerMessages", len(publishStorageDealsMessages))
			util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page="+fmt.Sprintf("%d", i)+"&method=PublishStorageDeals", filFoxMessagesList)
			publishStorageDealsMessages = append(publishStorageDealsMessages, filFoxMessagesList.Messages...)
		}
	} else if db_publish_storage_deals_messages_total_count != int64(totalPublishStorageDealsMessageCount) {
		minerMessagePagesCount := totalPublishStorageDealsMessageCount / 100
		fmt.Println("minerMessagePagesCount", minerMessagePagesCount)
		for i := 1; i <= minerMessagePagesCount; i++ {
			fmt.Println("page", i)
			fmt.Println("iterminerMessages", len(publishStorageDealsMessages))
			util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page="+fmt.Sprintf("%d", i)+"&method=PublishStorageDeals", filFoxMessagesList)
			publishStorageDealsMessages = append(publishStorageDealsMessages, filFoxMessagesList.Messages...)
		}
	}

	if db_publish_storage_deals_messages_total_count != int64(totalPublishStorageDealsMessageCount) {
		for _, psdm := range publishStorageDealsMessages {
			filFoxMessage := new(FilFoxPublishStorageDealsMessage)
			util.GetJson(FILFOX_MESSAGE+psdm.Cid, filFoxMessage)

			transactionType := "Deals Publish"
			burnFee := "0"
			if len(filFoxMessage.Transfers) >= 2 {
				burnFee = filFoxMessage.Transfers[1].Value
			}
			ids := filFoxMessage.DecodedReturnValue.IDs
			provider := filFoxMessage.DecodedParams.Deals[0].Proposal.Provider
			value := filFoxMessage.Value
			minerFee := filFoxMessage.Fee.MinerTip

			if value == "" {
				value = "0"
			}
			if minerFee == "" {
				minerFee = "0"
			}
			if burnFee == "" {
				burnFee = "0"
			}
			_, err := DB.Model(&model.Transaction{
				ID:              psdm.Cid,
				MinerID:         provider,
				Height:          int64(psdm.Height),
				Timestamp:       int64(psdm.Timestamp),
				TransactionType: transactionType,
				MethodName:      psdm.Method,
				Value:           value,
				MinerFee:        minerFee,
				BurnFee:         burnFee,
				From:            psdm.From,
				To:              psdm.To,
				ExitCode:        psdm.Receipt.ExitCode,
				Deals:           ids,
			}).Insert()
			if err != nil {
				fmt.Println("publishStorageDealsMessages insert err:", err)
			}
		}
	}

	_, err = DB.Model(&model.FilfoxMessagesCount{
		// ID:                      m,
		PublishStorageDealsMessagesTotalCount: int64(totalPublishStorageDealsMessageCount),
	}).
		// OnConflict("(id) DO UPDATE").
		// Set("title = EXCLUDED.title").
		Column("publish_storage_deals_messages_total_count").
		Where("id = 'psdm'").
		Update()
	if err != nil {
		fmt.Println("inserting/updating PublishStorageDealsMessagesTotalCount", err)
	}
}

func WithdrawBalanceMarketMessages(DB *pg.DB, node lens.API) {
	var FILFOX_MESSAGE string = "https://filfox.info/api/v1/message/"

	filFoxMessagesList := new(FilFoxMessagesList)
	util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page=0&method=WithdrawBalance%20(market)", filFoxMessagesList)

	withdrawBalanceMarketMessages := []struct {
		Cid       string `json:"cid"`
		Height    int    `json:"height"`
		Timestamp int    `json:"timestamp"`
		From      string `json:"from"`
		To        string `json:"to"`
		ToTag     struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag"`
		Value   string `json:"value"`
		Method  string `json:"method"`
		Receipt struct {
			ExitCode int `json:"exitCode"`
		} `json:"receipt"`
		FromTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"fromTag,omitempty"`
	}{}

	withdrawBalanceMarketMessages = append(withdrawBalanceMarketMessages, filFoxMessagesList.Messages...)

	var db_withdraw_balance_market_messages_total_count int64
	err := DB.Model((*model.FilfoxMessagesCount)(nil)).
		Column("withdraw_balance_market_messages_total_count").
		Where("id = 'wbmm'").
		Select(&db_withdraw_balance_market_messages_total_count)

	fmt.Println("db_withdraw_balance_market_messages_total_count", db_withdraw_balance_market_messages_total_count)

	totalWithdrawBalanceMarketMessageCount := filFoxMessagesList.TotalCount
	fmt.Println("totalWithdrawBalanceMarketMessageCount", totalWithdrawBalanceMarketMessageCount)

	var diff int64
	var pages int
	if err == nil && db_withdraw_balance_market_messages_total_count < int64(totalWithdrawBalanceMarketMessageCount) {
		diff = int64(totalWithdrawBalanceMarketMessageCount) - db_withdraw_balance_market_messages_total_count
		pages = int(diff) / 100
		fmt.Println("case1 diff", diff, "pages", pages)
		for i := 1; i <= pages; i++ {
			fmt.Println("page", i)
			fmt.Println("iterminerMessages", len(withdrawBalanceMarketMessages))
			util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page="+fmt.Sprintf("%d", i)+"&method=WithdrawBalance%20(market)", filFoxMessagesList)
			withdrawBalanceMarketMessages = append(withdrawBalanceMarketMessages, filFoxMessagesList.Messages...)
		}
	} else if db_withdraw_balance_market_messages_total_count != int64(totalWithdrawBalanceMarketMessageCount) {
		minerMessagePagesCount := totalWithdrawBalanceMarketMessageCount / 100
		fmt.Println("minerMessagePagesCount", minerMessagePagesCount)
		for i := 1; i <= minerMessagePagesCount; i++ {
			fmt.Println("page", i)
			fmt.Println("iterminerMessages", len(withdrawBalanceMarketMessages))
			util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page="+fmt.Sprintf("%d", i)+"&method=WithdrawBalance%20(market)", filFoxMessagesList)
			withdrawBalanceMarketMessages = append(withdrawBalanceMarketMessages, filFoxMessagesList.Messages...)
		}
	}

	if db_withdraw_balance_market_messages_total_count != int64(totalWithdrawBalanceMarketMessageCount) {
		for _, wbmm := range withdrawBalanceMarketMessages {
			filFoxMessage := new(FilFoxWithdrawBalanceMarketMessage)
			util.GetJson(FILFOX_MESSAGE+wbmm.Cid, filFoxMessage)

			burnFee := "0"
			if len(filFoxMessage.Transfers) >= 2 {
				burnFee = filFoxMessage.Transfers[1].Value
			}
			miner := filFoxMessage.DecodedParams.Address
			transactionType := "Transfer"
			value := filFoxMessage.DecodedParams.Amount
			minerFee := filFoxMessage.Fee.MinerTip

			if value == "" {
				value = "0"
			}
			if minerFee == "" {
				minerFee = "0"
			}
			if burnFee == "" {
				burnFee = "0"
			}
			_, err := DB.Model(&model.Transaction{
				ID:              wbmm.Cid,
				MinerID:         miner,
				Height:          int64(wbmm.Height),
				Timestamp:       int64(wbmm.Timestamp),
				TransactionType: transactionType,
				MethodName:      wbmm.Method,
				Value:           value,
				MinerFee:        minerFee,
				BurnFee:         burnFee,
				From:            wbmm.From,
				To:              wbmm.To,
				ExitCode:        wbmm.Receipt.ExitCode,
			}).Insert()
			if err != nil {
				fmt.Println("withdrawBalanceMarketMessages insert err:", err)
			}
		}
	}

	_, err = DB.Model(&model.FilfoxMessagesCount{
		// ID:                      m,
		WithdrawBalanceMarketMessagesTotalCount: int64(totalWithdrawBalanceMarketMessageCount),
	}).
		// OnConflict("(id) DO UPDATE").
		// Set("title = EXCLUDED.title").
		Column("withdraw_balance_market_messages_total_count").
		Where("id = 'wbmm'").
		Update()
	if err != nil {
		fmt.Println("inserting/updating WithdrawBalanceMarketMessagesTotalCount", err)
	}
}

func AddBalanceMessages(DB *pg.DB, node lens.API) {
	var FILFOX_MESSAGE string = "https://filfox.info/api/v1/message/"

	filFoxMessagesList := new(FilFoxMessagesList)
	util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page=0&method=AddBalance", filFoxMessagesList)

	addBalanceMessages := []struct {
		Cid       string `json:"cid"`
		Height    int    `json:"height"`
		Timestamp int    `json:"timestamp"`
		From      string `json:"from"`
		To        string `json:"to"`
		ToTag     struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag"`
		Value   string `json:"value"`
		Method  string `json:"method"`
		Receipt struct {
			ExitCode int `json:"exitCode"`
		} `json:"receipt"`
		FromTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"fromTag,omitempty"`
	}{}

	addBalanceMessages = append(addBalanceMessages, filFoxMessagesList.Messages...)

	var db_add_balance_messages_total_count int64
	err := DB.Model((*model.FilfoxMessagesCount)(nil)).
		Column("add_balance_messages_total_count").
		Where("id = 'abm'").
		Select(&db_add_balance_messages_total_count)

	fmt.Println("db_add_balance_messages_total_count", db_add_balance_messages_total_count)

	totalAddBalanceMessageCount := filFoxMessagesList.TotalCount
	fmt.Println("totalAddBalanceMessageCount", totalAddBalanceMessageCount)

	var diff int64
	var pages int
	if err == nil && db_add_balance_messages_total_count < int64(totalAddBalanceMessageCount) {
		diff = int64(totalAddBalanceMessageCount) - db_add_balance_messages_total_count
		pages = int(diff) / 100
		fmt.Println("case1 diff", diff, "pages", pages)
		for i := 1; i <= pages; i++ {
			fmt.Println("page", i)
			fmt.Println("iterminerMessages", len(addBalanceMessages))
			util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page="+fmt.Sprintf("%d", i)+"&method=AddBalance", filFoxMessagesList)
			addBalanceMessages = append(addBalanceMessages, filFoxMessagesList.Messages...)
		}
	} else if db_add_balance_messages_total_count != int64(totalAddBalanceMessageCount) {
		minerMessagePagesCount := totalAddBalanceMessageCount / 100
		fmt.Println("minerMessagePagesCount", minerMessagePagesCount)
		for i := 1; i <= minerMessagePagesCount; i++ {
			fmt.Println("page", i)
			fmt.Println("iterminerMessages", len(addBalanceMessages))
			util.GetJson(FILFOX_MESSAGE+"list?pageSize=100&page="+fmt.Sprintf("%d", i)+"&method=AddBalance", filFoxMessagesList)
			addBalanceMessages = append(addBalanceMessages, filFoxMessagesList.Messages...)
		}
	}

	if db_add_balance_messages_total_count != int64(totalAddBalanceMessageCount) {
		for _, abm := range addBalanceMessages {
			filFoxAddBalanceMessage := new(FilFoxAddBalanceMessage)
			util.GetJson(FILFOX_MESSAGE+abm.Cid, filFoxAddBalanceMessage)

			burnFee := "0"
			if len(filFoxAddBalanceMessage.Transfers) >= 2 {
				burnFee = filFoxAddBalanceMessage.Transfers[1].Value
			}
			provider := filFoxAddBalanceMessage.DecodedParams
			transactionType := "Collateral Deposit"
			value := filFoxAddBalanceMessage.Value
			minerFee := filFoxAddBalanceMessage.Fee.MinerTip

			if value == "" {
				value = "0"
			}
			if minerFee == "" {
				minerFee = "0"
			}
			if burnFee == "" {
				burnFee = "0"
			}
			_, err := DB.Model(&model.Transaction{
				ID:              abm.Cid,
				MinerID:         provider,
				Height:          int64(abm.Height),
				Timestamp:       int64(abm.Timestamp),
				TransactionType: transactionType,
				MethodName:      abm.Method,
				Value:           value,
				MinerFee:        minerFee,
				BurnFee:         burnFee,
				From:            abm.From,
				To:              abm.To,
				ExitCode:        abm.Receipt.ExitCode,
			}).Insert()
			if err != nil {
				fmt.Println("addBalanceMessages insert err:", err)
			}
		}
	}

	_, err = DB.Model(&model.FilfoxMessagesCount{
		// ID:                      m,
		AddBalanceMessagesTotalCount: int64(totalAddBalanceMessageCount),
	}).
		// OnConflict("(id) DO UPDATE").
		// Set("title = EXCLUDED.title").
		Column("add_balance_messages_total_count").
		Where("id = 'abm'").
		Update()
	if err != nil {
		fmt.Println("inserting/updating AddBalanceMessagesTotalCount", err)
	}
}

type FilFoxStatsPower []struct {
	Height               int    `json:"height"`
	Timestamp            int    `json:"timestamp"`
	RawBytePower         string `json:"rawBytePower"`
	QualityAdjPower      string `json:"qualityAdjPower"`
	RawBytePowerDelta    string `json:"rawBytePowerDelta"`
	QualityAdjPowerDelta string `json:"qualityAdjPowerDelta"`
}

type FilFoxMinerBlocks struct {
	TotalCount int `json:"totalCount"`
	Blocks     []struct {
		Cid          string `json:"cid"`
		Height       int    `json:"height"`
		Timestamp    int    `json:"timestamp"`
		Size         int    `json:"size"`
		WinCount     int    `json:"winCount"`
		Reward       string `json:"reward"`
		Penalty      string `json:"penalty"`
		MessageCount int    `json:"messageCount"`
	} `json:"blocks"`
}

type FilFoxMinerTransfers struct {
	TotalCount int `json:"totalCount"`
	Transfers  []struct {
		Height    int    `json:"height"`
		Timestamp int    `json:"timestamp"`
		From      string `json:"from"`
		FromTag   struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"fromTag"`
		To    string `json:"to"`
		ToTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"transfers"`
	Types []string `json:"types"`
}

type FilFoxMinerMessages struct {
	TotalCount int `json:"totalCount"`
	Messages   []struct {
		Cid       string `json:"cid"`
		Height    int    `json:"height"`
		Timestamp int    `json:"timestamp"`
		From      string `json:"from"`
		To        string `json:"to"`
		Nonce     int    `json:"nonce"`
		Value     string `json:"value"`
		Method    string `json:"method"`
		Receipt   struct {
			ExitCode int `json:"exitCode"`
		} `json:"receipt"`
	} `json:"messages"`
	Methods []string `json:"methods"`
}

type FilFoxAddBalanceMessage struct {
	Cid           string   `json:"cid"`
	Height        int      `json:"height"`
	Timestamp     int      `json:"timestamp"`
	Confirmations int      `json:"confirmations"`
	Blocks        []string `json:"blocks"`
	Version       int      `json:"version"`
	From          string   `json:"from"`
	FromID        string   `json:"fromId"`
	FromActor     string   `json:"fromActor"`
	To            string   `json:"to"`
	ToID          string   `json:"toId"`
	ToActor       string   `json:"toActor"`
	ToTag         struct {
		Name   string `json:"name"`
		Signed bool   `json:"signed"`
	} `json:"toTag"`
	Nonce        int    `json:"nonce"`
	Value        string `json:"value"`
	GasLimit     int    `json:"gasLimit"`
	GasFeeCap    string `json:"gasFeeCap"`
	GasPremium   string `json:"gasPremium"`
	Method       string `json:"method"`
	MethodNumber int    `json:"methodNumber"`
	Params       string `json:"params"`
	Receipt      struct {
		ExitCode int    `json:"exitCode"`
		Return   string `json:"return"`
		GasUsed  int    `json:"gasUsed"`
	} `json:"receipt"`
	DecodedParams      string      `json:"decodedParams"`
	DecodedReturnValue interface{} `json:"decodedReturnValue"`
	Size               int         `json:"size"`
	Error              string      `json:"error"`
	BaseFee            string      `json:"baseFee"`
	Fee                struct {
		BaseFeeBurn        string `json:"baseFeeBurn"`
		OverEstimationBurn string `json:"overEstimationBurn"`
		MinerPenalty       string `json:"minerPenalty"`
		MinerTip           string `json:"minerTip"`
		Refund             string `json:"refund"`
	} `json:"fee"`
	Subcalls []struct {
		From         string `json:"from"`
		FromID       string `json:"fromId"`
		FromActor    string `json:"fromActor"`
		To           string `json:"to"`
		ToID         string `json:"toId"`
		ToActor      string `json:"toActor"`
		Value        string `json:"value"`
		Method       string `json:"method"`
		MethodNumber int    `json:"methodNumber"`
		Params       string `json:"params"`
		Receipt      struct {
			ExitCode int    `json:"exitCode"`
			Return   string `json:"return"`
			GasUsed  int    `json:"gasUsed"`
		} `json:"receipt"`
		DecodedParams      interface{} `json:"decodedParams"`
		DecodedReturnValue struct {
			Owner        string   `json:"Owner"`
			Worker       string   `json:"Worker"`
			ControlAddrs []string `json:"ControlAddrs"`
		} `json:"decodedReturnValue"`
		Error    string        `json:"error"`
		Subcalls []interface{} `json:"subcalls"`
	} `json:"subcalls"`
	Transfers []struct {
		From   string `json:"from"`
		FromID string `json:"fromId"`
		To     string `json:"to"`
		ToID   string `json:"toId"`
		ToTag  struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"transfers"`
}

type FilFoxWithdrawBalanceMarketMessage struct {
	Cid           string   `json:"cid"`
	Height        int      `json:"height"`
	Timestamp     int      `json:"timestamp"`
	Confirmations int      `json:"confirmations"`
	Blocks        []string `json:"blocks"`
	Version       int      `json:"version"`
	From          string   `json:"from"`
	FromID        string   `json:"fromId"`
	FromActor     string   `json:"fromActor"`
	FromTag       struct {
		Name   string `json:"name"`
		Signed bool   `json:"signed"`
	} `json:"fromTag"`
	To      string `json:"to"`
	ToID    string `json:"toId"`
	ToActor string `json:"toActor"`
	ToTag   struct {
		Name   string `json:"name"`
		Signed bool   `json:"signed"`
	} `json:"toTag"`
	Nonce        int    `json:"nonce"`
	Value        string `json:"value"`
	GasLimit     int    `json:"gasLimit"`
	GasFeeCap    string `json:"gasFeeCap"`
	GasPremium   string `json:"gasPremium"`
	Method       string `json:"method"`
	MethodNumber int    `json:"methodNumber"`
	Params       string `json:"params"`
	Receipt      struct {
		ExitCode int    `json:"exitCode"`
		Return   string `json:"return"`
		GasUsed  int    `json:"gasUsed"`
	} `json:"receipt"`
	DecodedParams struct {
		Address string `json:"Address"`
		Amount  string `json:"Amount"`
	} `json:"decodedParams"`
	DecodedReturnValue interface{} `json:"decodedReturnValue"`
	Size               int         `json:"size"`
	Error              string      `json:"error"`
	BaseFee            string      `json:"baseFee"`
	Fee                struct {
		BaseFeeBurn        string `json:"baseFeeBurn"`
		OverEstimationBurn string `json:"overEstimationBurn"`
		MinerPenalty       string `json:"minerPenalty"`
		MinerTip           string `json:"minerTip"`
		Refund             string `json:"refund"`
	} `json:"fee"`
	Subcalls []struct {
		From         string `json:"from"`
		FromID       string `json:"fromId"`
		FromActor    string `json:"fromActor"`
		To           string `json:"to"`
		ToID         string `json:"toId"`
		ToActor      string `json:"toActor"`
		Value        string `json:"value"`
		Method       string `json:"method"`
		MethodNumber int    `json:"methodNumber"`
		Params       string `json:"params"`
		Receipt      struct {
			ExitCode int    `json:"exitCode"`
			Return   string `json:"return"`
			GasUsed  int    `json:"gasUsed"`
		} `json:"receipt"`
		DecodedParams      interface{}   `json:"decodedParams"`
		DecodedReturnValue interface{}   `json:"decodedReturnValue"`
		Error              string        `json:"error"`
		Subcalls           []interface{} `json:"subcalls"`
	} `json:"subcalls"`
	Transfers []struct {
		From    string `json:"from"`
		FromID  string `json:"fromId"`
		FromTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"fromTag"`
		To    string `json:"to"`
		ToID  string `json:"toId"`
		Value string `json:"value"`
		Type  string `json:"type"`
		ToTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag,omitempty"`
	} `json:"transfers"`
}

type FilFoxPublishStorageDealsMessage struct {
	Cid           string   `json:"cid"`
	Height        int      `json:"height"`
	Timestamp     int      `json:"timestamp"`
	Confirmations int      `json:"confirmations"`
	Blocks        []string `json:"blocks"`
	Version       int      `json:"version"`
	From          string   `json:"from"`
	FromID        string   `json:"fromId"`
	FromActor     string   `json:"fromActor"`
	To            string   `json:"to"`
	ToID          string   `json:"toId"`
	ToActor       string   `json:"toActor"`
	ToTag         struct {
		Name   string `json:"name"`
		Signed bool   `json:"signed"`
	} `json:"toTag"`
	Nonce        int    `json:"nonce"`
	Value        string `json:"value"`
	GasLimit     int    `json:"gasLimit"`
	GasFeeCap    string `json:"gasFeeCap"`
	GasPremium   string `json:"gasPremium"`
	Method       string `json:"method"`
	MethodNumber int    `json:"methodNumber"`
	Params       string `json:"params"`
	Receipt      struct {
		ExitCode int    `json:"exitCode"`
		Return   string `json:"return"`
		GasUsed  int    `json:"gasUsed"`
	} `json:"receipt"`
	DecodedParams struct {
		Deals []struct {
			Proposal struct {
				PieceCID             string `json:"PieceCID"`
				PieceSize            int64  `json:"PieceSize"`
				VerifiedDeal         bool   `json:"VerifiedDeal"`
				Client               string `json:"Client"`
				Provider             string `json:"Provider"`
				Label                string `json:"Label"`
				StartEpoch           int    `json:"StartEpoch"`
				EndEpoch             int    `json:"endEpoch"`
				StoragePricePerEpoch string `json:"StoragePricePerEpoch"`
				ProviderCollateral   string `json:"ProviderCollateral"`
				ClientCollateral     string `json:"ClientCollateral"`
			} `json:"Proposal"`
			ClientSignature struct {
				Type      int    `json:"Type"`
				Signature string `json:"Signature"`
			} `json:"ClientSignature"`
		} `json:"Deals"`
	} `json:"decodedParams"`
	DecodedReturnValue struct {
		IDs []int `json:"IDs"`
	} `json:"decodedReturnValue"`
	Size    int    `json:"size"`
	Error   string `json:"error"`
	BaseFee string `json:"baseFee"`
	Fee     struct {
		BaseFeeBurn        string `json:"baseFeeBurn"`
		OverEstimationBurn string `json:"overEstimationBurn"`
		MinerPenalty       string `json:"minerPenalty"`
		MinerTip           string `json:"minerTip"`
		Refund             string `json:"refund"`
	} `json:"fee"`
	Subcalls []struct {
		From         string `json:"from"`
		FromID       string `json:"fromId"`
		FromActor    string `json:"fromActor"`
		To           string `json:"to"`
		ToID         string `json:"toId"`
		ToActor      string `json:"toActor"`
		Value        string `json:"value"`
		Method       string `json:"method"`
		MethodNumber int    `json:"methodNumber"`
		Params       string `json:"params"`
		Receipt      struct {
			ExitCode int    `json:"exitCode"`
			Return   string `json:"return"`
			GasUsed  int    `json:"gasUsed"`
		} `json:"receipt"`
		DecodedParams      interface{} `json:"decodedParams"`
		DecodedReturnValue struct {
			Owner                   string   `json:"Owner"`
			Worker                  string   `json:"Worker"`
			ControlAddrs            []string `json:"ControlAddrs"`
			RawBytePower            string   `json:"RawBytePower"`
			QualityAdjPower         string   `json:"QualityAdjPower"`
			PledgeCollateral        string   `json:"PledgeCollateral"`
			QualityAdjPowerSmoothed struct {
				PositionEstimate string `json:"PositionEstimate"`
				VelocityEstimate string `json:"VelocityEstimate"`
			} `json:"QualityAdjPowerSmoothed"`
			ThisEpochRewardSmoothed struct {
				PositionEstimate string `json:"PositionEstimate"`
				VelocityEstimate string `json:"VelocityEstimate"`
			} `json:"ThisEpochRewardSmoothed"`
			ThisEpochBaselinePower string `json:"ThisEpochBaselinePower"`
		} `json:"decodedReturnValue,omitempty"`
		Error    string        `json:"error"`
		Subcalls []interface{} `json:"subcalls"`
	} `json:"subcalls"`
	Transfers []struct {
		From   string `json:"from"`
		FromID string `json:"fromId"`
		To     string `json:"to"`
		ToID   string `json:"toId"`
		ToTag  struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"transfers"`
}

type FilFoxMessage struct {
	Cid           string   `json:"cid"`
	Height        int      `json:"height"`
	Timestamp     int      `json:"timestamp"`
	Confirmations int      `json:"confirmations"`
	Blocks        []string `json:"blocks"`
	Version       int      `json:"version"`
	From          string   `json:"from"`
	FromID        string   `json:"fromId"`
	FromActor     string   `json:"fromActor"`
	FromTag       struct {
		Name   string `json:"name"`
		Signed bool   `json:"signed"`
	} `json:"fromTag"`
	To      string `json:"to"`
	ToID    string `json:"toId"`
	ToActor string `json:"toActor"`
	ToTag   struct {
		Name   string `json:"name"`
		Signed bool   `json:"signed"`
	} `json:"toTag"`
	Nonce        int    `json:"nonce"`
	Value        string `json:"value"`
	GasLimit     int    `json:"gasLimit"`
	GasFeeCap    string `json:"gasFeeCap"`
	GasPremium   string `json:"gasPremium"`
	Method       string `json:"method"`
	MethodNumber int    `json:"methodNumber"`
	Params       string `json:"params"`
	Receipt      struct {
		ExitCode int    `json:"exitCode"`
		Return   string `json:"return"`
		GasUsed  int    `json:"gasUsed"`
	} `json:"receipt"`
	DecodedParams struct {
		Address         string `json:"Address"`
		Amount          string `json:"Amount"`
		AmountRequested string `json:"AmountRequested"`
		Deadline        int    `json:"Deadline"`
		Deals           []struct {
			Proposal struct {
				PieceCID             string `json:"PieceCID"`
				PieceSize            int64  `json:"PieceSize"`
				VerifiedDeal         bool   `json:"VerifiedDeal"`
				Client               string `json:"Client"`
				Provider             string `json:"Provider"`
				Label                string `json:"Label"`
				StartEpoch           int    `json:"StartEpoch"`
				EndEpoch             int    `json:"endEpoch"`
				StoragePricePerEpoch string `json:"StoragePricePerEpoch"`
				ProviderCollateral   string `json:"ProviderCollateral"`
				ClientCollateral     string `json:"ClientCollateral"`
			} `json:"Proposal"`
			ClientSignature struct {
				Type      int    `json:"Type"`
				Signature string `json:"Signature"`
			} `json:"ClientSignature"`
		} `json:"Deals"`
		PoStIndex  int `json:"PoStIndex"`
		Partitions []struct {
			Index   int    `json:"Index"`
			Skipped string `json:"Skipped"`
		} `json:"Partitions"`
		Proofs []struct {
			PoStProof  int    `json:"PoStProof"`
			ProofBytes string `json:"ProofBytes"`
		} `json:"Proofs"`
		ChainCommitEpoch int    `json:"ChainCommitEpoch"`
		ChainCommitRand  string `json:"ChainCommitRand"`
	} `json:"decodedParams"`
	DecodedReturnValue struct {
		IDs []int `json:"IDs"`
	} `json:"decodedReturnValue"`
	Size    int    `json:"size"`
	Error   string `json:"error"`
	BaseFee string `json:"baseFee"`
	Fee     struct {
		BaseFeeBurn        string `json:"baseFeeBurn"`
		OverEstimationBurn string `json:"overEstimationBurn"`
		MinerPenalty       string `json:"minerPenalty"`
		MinerTip           string `json:"minerTip"`
		Refund             string `json:"refund"`
	} `json:"fee"`
	Subcalls []struct {
		From         string `json:"from"`
		FromID       string `json:"fromId"`
		FromActor    string `json:"fromActor"`
		To           string `json:"to"`
		ToID         string `json:"toId"`
		ToActor      string `json:"toActor"`
		Value        string `json:"value"`
		Method       string `json:"method"`
		MethodNumber int    `json:"methodNumber"`
		Params       string `json:"params"`
		Receipt      struct {
			ExitCode int    `json:"exitCode"`
			Return   string `json:"return"`
			GasUsed  int    `json:"gasUsed"`
		} `json:"receipt"`
		DecodedParams      interface{}   `json:"decodedParams"`
		Error              string        `json:"error"`
		Subcalls           []interface{} `json:"subcalls"`
		DecodedReturnValue struct {
			Owner                   string   `json:"Owner"`
			Worker                  string   `json:"Worker"`
			ControlAddrs            []string `json:"ControlAddrs"`
			RawBytePower            string   `json:"RawBytePower"`
			QualityAdjPower         string   `json:"QualityAdjPower"`
			PledgeCollateral        string   `json:"PledgeCollateral"`
			QualityAdjPowerSmoothed struct {
				PositionEstimate string `json:"PositionEstimate"`
				VelocityEstimate string `json:"VelocityEstimate"`
			} `json:"QualityAdjPowerSmoothed"`
			ThisEpochRewardSmoothed struct {
				PositionEstimate string `json:"PositionEstimate"`
				VelocityEstimate string `json:"VelocityEstimate"`
			} `json:"ThisEpochRewardSmoothed"`
			ThisEpochBaselinePower string `json:"ThisEpochBaselinePower"`
		} `json:"decodedReturnValue,omitempty"`
	} `json:"subcalls"`
	Transfers []struct {
		From   string `json:"from"`
		FromID string `json:"fromId"`
		To     string `json:"to"`
		ToID   string `json:"toId"`
		ToTag  struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag,omitempty"`
		Value   string `json:"value"`
		Type    string `json:"type"`
		FromTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"fromTag,omitempty"`
	} `json:"transfers"`
}

type FilFoxMessagesList struct {
	TotalCount int `json:"totalCount"`
	Messages   []struct {
		Cid       string `json:"cid"`
		Height    int    `json:"height"`
		Timestamp int    `json:"timestamp"`
		From      string `json:"from"`
		To        string `json:"to"`
		ToTag     struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"toTag"`
		Value   string `json:"value"`
		Method  string `json:"method"`
		Receipt struct {
			ExitCode int `json:"exitCode"`
		} `json:"receipt"`
		FromTag struct {
			Name   string `json:"name"`
			Signed bool   `json:"signed"`
		} `json:"fromTag,omitempty"`
	} `json:"messages"`
	Methods []string `json:"methods"`
}

func GetMessageAttributes(node lens.API, filfoxMessage FilFoxMessage) (string, string, string, string, string, []int) {
	fmt.Println("inside GetMessageAttributes", filfoxMessage)
	switch filfoxMessage.Method {
	case "PreCommitSector", "ProveCommitSector":
		burnFee := "0"
		if len(filfoxMessage.Transfers) >= 2 {
			burnFee = filfoxMessage.Transfers[1].Value
		}
		return "Collateral Deposit", filfoxMessage.Value, filfoxMessage.Fee.MinerTip, burnFee, "", []int{}
	case "ReportConsensusFault", "DisputeWindowedPoSt":
		transfer, _ := strconv.ParseInt(filfoxMessage.Transfers[2].Value, 10, 64)
		burn, _ := strconv.ParseInt(filfoxMessage.Transfers[3].Value, 10, 64)
		amt := -1 * (transfer + burn)
		burnFee := "0"
		if len(filfoxMessage.Transfers) >= 2 {
			burnFee = filfoxMessage.Transfers[1].Value
		}
		return "Penalty", fmt.Sprintf("%v", amt), filfoxMessage.Fee.MinerTip, burnFee, "", []int{}
	case "TerminateSectors":
		burnFee := "0"
		if len(filfoxMessage.Transfers) >= 2 {
			burnFee = filfoxMessage.Transfers[1].Value
		}
		return "Penalty", filfoxMessage.Transfers[2].Value, filfoxMessage.Fee.MinerTip, burnFee, "", []int{}
	case "RepayDebt":
		burnFee := "0"
		if len(filfoxMessage.Transfers) >= 2 {
			burnFee = filfoxMessage.Transfers[1].Value
		}
		return "Penalty", filfoxMessage.Value, filfoxMessage.Fee.MinerTip, burnFee, "", []int{}
	case "WithdrawBalance (miner)":
		burnFee := "0"
		if len(filfoxMessage.Transfers) >= 2 {
			burnFee = filfoxMessage.Transfers[1].Value
		}
		return "Transfer", filfoxMessage.DecodedParams.AmountRequested, filfoxMessage.Fee.MinerTip, burnFee, "", []int{}
	// case "WithdrawBalance (market)":
	// 	burnFee := "0"
	// 	if len(filfoxMessage.Transfers) >= 2 {
	// 		burnFee = filfoxMessage.Transfers[1].Value
	// 	}
	// 	miner := filfoxMessage.DecodedParams.Address
	// 	return "Transfer", filfoxMessage.DecodedParams.Amount, filfoxMessage.Fee.MinerTip, burnFee, miner, []int{}

	// case "PublishStorageDeals":
	// 	burnFee := "0"
	// 	if len(filfoxMessage.Transfers) >= 2 {
	// 		burnFee = filfoxMessage.Transfers[1].Value
	// 	}
	// 	ids := filfoxMessage.DecodedReturnValue.IDs
	// 	provider := filfoxMessage.DecodedParams.Deals[0].Proposal.Provider
	// 	return "Deals Publish", filfoxMessage.Value, filfoxMessage.Fee.MinerTip, burnFee, provider, ids

	// case "AddBalance":
	// 	burnFee := "0"
	// 	if len(filfoxMessage.Transfers) >= 2 {
	// 		burnFee = filfoxMessage.Transfers[1].Value
	// 	}
	// 	provider := filfoxMessage.DecodedParams

	// return "Collateral Deposit", filfoxMessage.Value, filfoxMessage.Fee.MinerTip, burnFee, "", []int{}
	// case "SubmitWindowedPoSt", "ChangeWorkerAddress", "ChangePeerID", "ExtendSectorExpiration",
	// 	"DeclareFaults", "DeclareFaultsRecovered", "ChangeMultiaddrs", "CompactSectorNumbers",
	// 	"ConfirmUpdateWorkerKey", "ChangeOwnerAddress", "CreateMiner":
	default:
		basefeeburn, _ := strconv.ParseInt(filfoxMessage.Fee.BaseFeeBurn, 10, 64)
		overestimationburn, _ := strconv.ParseInt(filfoxMessage.Fee.OverEstimationBurn, 10, 64)
		burnfee := basefeeburn + overestimationburn
		return "Others", filfoxMessage.Value, filfoxMessage.Fee.MinerTip, fmt.Sprintf("%d", burnfee), "", []int{}
	}
	// return transationType, value, minerFee, burnFee
}
