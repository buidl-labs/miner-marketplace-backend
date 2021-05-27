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
	StorageAskPrice      float64
	VerifiedAskPrice     float64
	RetrievalAskPrice    float64
	ReputationScore      int `pg:",notnull"`
	TransparencyScore    int `pg:",notnull"`
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
