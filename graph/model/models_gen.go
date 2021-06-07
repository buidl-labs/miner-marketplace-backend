// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type DataTransferMechanism struct {
	Online  bool `json:"online"`
	Offline bool `json:"offline"`
}

type Location struct {
	Region  string `json:"region"`
	Country string `json:"country"`
}

type Miner struct {
	ID                   string        `json:"id"`
	Claimed              bool          `json:"claimed"`
	PersonalInfo         *PersonalInfo `json:"personalInfo"`
	Worker               *Worker       `json:"worker"`
	Owner                *Owner        `json:"owner"`
	Location             *Location     `json:"location"`
	QualityAdjustedPower string        `json:"qualityAdjustedPower"`
	Service              *Service      `json:"service"`
	Pricing              *Pricing      `json:"pricing"`
	ReputationScore      int           `json:"reputationScore"`
	TransparencyScore    int           `json:"transparencyScore"`
}

type NetworkStats struct {
	ActiveMinersCount      int    `json:"activeMinersCount"`
	NetworkStorageCapacity string `json:"networkStorageCapacity"`
	DataStored             string `json:"dataStored"`
}

type Owner struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Miner   *Miner `json:"miner"`
}

type PersonalInfo struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Twitter string `json:"twitter"`
	Slack   string `json:"slack"`
}

type Pricing struct {
	StorageAskPrice   string `json:"storageAskPrice"`
	VerifiedAskPrice  string `json:"verifiedAskPrice"`
	RetrievalAskPrice string `json:"retrievalAskPrice"`
}

type ProfileClaimInput struct {
	MinerID       string `json:"minerID"`
	LedgerAddress string `json:"ledgerAddress"`
}

type ProfileSettingsInput struct {
	MinerID           string `json:"minerID"`
	LedgerAddress     string `json:"ledgerAddress"`
	Name              string `json:"name"`
	Bio               string `json:"bio"`
	Email             string `json:"email"`
	Website           string `json:"website"`
	Twitter           string `json:"twitter"`
	Slack             string `json:"slack"`
	Region            string `json:"region"`
	Country           string `json:"country"`
	Storage           bool   `json:"storage"`
	Retrieval         bool   `json:"retrieval"`
	Repair            bool   `json:"repair"`
	Online            bool   `json:"online"`
	Offline           bool   `json:"offline"`
	StorageAskPrice   string `json:"storageAskPrice"`
	VerifiedAskPrice  string `json:"verifiedAskPrice"`
	RetrievalAskPrice string `json:"retrievalAskPrice"`
}

type Service struct {
	ServiceTypes          *ServiceTypes          `json:"serviceTypes"`
	DataTransferMechanism *DataTransferMechanism `json:"dataTransferMechanism"`
}

type ServiceTypes struct {
	Storage   bool `json:"storage"`
	Retrieval bool `json:"retrieval"`
	Repair    bool `json:"repair"`
}

type Worker struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Miner   *Miner `json:"miner"`
}
