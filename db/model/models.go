package model

type JoinedMiner struct {
	tableName            struct{} `pg:"miners"`
	ID                   string   `pg:",pk,notnull"`
	Claimed              bool     `pg:",notnull"`
	Onboarded            bool
	Region               string
	Country              string
	WorkerID             string
	WorkerAddress        string
	OwnerID              string
	OwnerAddress         string
	QualityAdjustedPower string `pg:",notnull"`
	StorageAskPrice      string
	VerifiedAskPrice     string
	RetrievalAskPrice    string
	ReputationScore      int    `pg:",notnull,use_zero"`
	TransparencyScore    int    `pg:",notnull,use_zero"`
	Storage              bool   `pg:",notnull"`
	Retrieval            bool   `pg:",notnull"`
	Repair               bool   `pg:",notnull"`
	DataTransferOnline   bool   `pg:",notnull"`
	DataTransferOffline  bool   `pg:",notnull"`
	Name                 string `pg:",notnull"`
	Bio                  string `pg:",notnull"`
	Email                string `pg:",notnull"`
	Website              string `pg:",notnull"`
	Twitter              string `pg:",notnull"`
	Slack                string `pg:",notnull"`
}

type Miner struct {
	ID                   string `pg:",pk,notnull"`
	Claimed              bool   `pg:",notnull"`
	Onboarded            bool
	Region               string
	Country              string
	WorkerID             string
	WorkerAddress        string
	OwnerID              string
	OwnerAddress         string
	QualityAdjustedPower string `pg:",notnull"`
	StorageAskPrice      string
	VerifiedAskPrice     string
	RetrievalAskPrice    string
	ReputationScore      int `pg:",notnull,use_zero"`
	TransparencyScore    int `pg:",notnull,use_zero"`
}

type MinerPersonalInfo struct {
	ID      string `pg:",pk,notnull"`
	Name    string `pg:",notnull"`
	Bio     string `pg:",notnull"`
	Email   string `pg:",notnull"`
	Website string `pg:",notnull"`
	Twitter string `pg:",notnull"`
	Slack   string `pg:",notnull"`
}

type MinerService struct {
	ID                  string `pg:",pk,notnull"`
	Storage             bool   `pg:",notnull"`
	Retrieval           bool   `pg:",notnull"`
	Repair              bool   `pg:",notnull"`
	DataTransferOnline  bool   `pg:",notnull"`
	DataTransferOffline bool   `pg:",notnull"`
}

type Transaction struct {
	ID              string `pg:",pk,notnull"`
	MinerID         string `pg:",notnull"`
	Height          int64  `pg:",notnull"`
	Timestamp       int64  `pg:",use_zero"`
	TransactionType string `pg:",notnull"`
	MethodName      string `pg:",notnull"`
	Value           string `pg:",notnull,use_zero"`
	MinerFee        string `pg:",notnull,use_zero"`
	BurnFee         string `pg:",notnull,use_zero"`
	From            string `pg:",notnull"`
	To              string `pg:",notnull"`
	ExitCode        int    `pg:",notnull,use_zero"`
	Deals           []int  `pg:",array"`
}

type MinerStorageDealStats struct {
	ID              string `pg:",pk,notnull"`
	AveragePrice    string `pg:",use_zero"`
	DataStored      string `pg:",use_zero"`
	FaultTerminated int64  `pg:",use_zero"`
	NoPenalties     int64  `pg:",use_zero"`
	Slashed         int64  `pg:",use_zero"`
	SuccessRate     string `pg:",use_zero"`
	Terminated      int64  `pg:",use_zero"`
	Total           int64  `pg:",use_zero"`
}

type MarketDealProposal struct {
	ID                   uint64 `pg:",pk,notnull,use_zero"`
	Height               int64  `pg:",use_zero"`
	Timestamp            int64  `pg:",use_zero"`
	PieceCID             string
	PieceSize            uint64 `pg:",notnull,use_zero"`
	VerifiedDeal         bool   `pg:",notnull,use_zero"`
	Provider             string `pg:",notnull"`
	Client               string `pg:",notnull"`
	Label                string
	StartEpoch           int64  `pg:",notnull,use_zero"`
	EndEpoch             int64  `pg:",notnull,use_zero"`
	StartTimestamp       int64  `pg:",use_zero"`
	EndTimestamp         int64  `pg:",use_zero"`
	StoragePrice         string `pg:",notnull,use_zero"`
	StoragePricePerEpoch string `pg:",use_zero"`
	ProviderCollateral   string `pg:",use_zero"`
	ClientCollateral     string `pg:",use_zero"`
}

type FilfoxMinerMessagesCount struct {
	ID                             string `pg:",pk,notnull"`
	MinerMessagesTotalCount        int64  `pg:",use_zero"`
	MinerTransfersRewardTotalCount int64  `pg:",use_zero"`
	TillHeight                     int64  `pg:","`
}

type FilfoxMessagesCount struct {
	ID                                      string `pg:",pk,notnull"`
	PublishStorageDealsMessagesTotalCount   int64  `pg:",use_zero"`
	WithdrawBalanceMarketMessagesTotalCount int64  `pg:",use_zero"`
	AddBalanceMessagesTotalCount            int64  `pg:",use_zero"`
	DealsTotalCount                         int64  `pg:",use_zero"`
}
