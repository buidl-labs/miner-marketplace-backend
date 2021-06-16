package model

type Miner struct {
	ID                   string `pg:",pk,notnull"`
	Claimed              bool   `pg:",notnull"`
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

type FilfoxMinerMessagesCount struct {
	ID                             string `pg:",pk,notnull"`
	MinerMessagesTotalCount        int64  `pg:",use_zero"`
	MinerTransfersRewardTotalCount int64  `pg:",use_zero"`
}

type FilfoxMessagesCount struct {
	ID                                      string `pg:",pk,notnull"`
	PublishStorageDealsMessagesTotalCount   int64  `pg:",use_zero"`
	WithdrawBalanceMarketMessagesTotalCount int64  `pg:",use_zero"`
	AddBalanceMessagesTotalCount            int64  `pg:",use_zero"`
}
