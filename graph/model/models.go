package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Contact struct {
	Miner   *Miner `json:"miner"`
	Email   string `json:"email"`
	Slack   string `json:"slack"`
	Website string `json:"website"`
	Twitter string `json:"twitter"`
}

type Deadline struct {
	ID            string `json:"id"`
	DeadlineIndex int64  `json:"deadlineIndex"`
	PeriodStart   int64  `json:"periodStart"`
	Open          int64  `json:"open"`
	Close         int64  `json:"close"`
	Challenge     int64  `json:"challenge"`
	FaultCutoff   int64  `json:"faultCutoff"`
}

type DeadlineOrderByInput struct {
	Deadline *Sort `json:"deadline"`
}

type Expenditure struct {
	NetworkFee string `json:"networkFee"`
	Penalty    string `json:"penalty"`
}

type Fault struct {
	Type      *FaultType `json:"type"`
	Penalty   *Penalty   `json:"penalty"`
	Height    int64      `json:"height"`
	Timestamp time.Time  `json:"timestamp"`
}

type FinanceMetrics struct {
	// ID          string       `json:"id"`
	// Income      *Income      `json:"income"`
	// Expenditure *Expenditure `json:"expenditure"`
	// Funds       *Funds       `json:"funds"`
	Miner                 *Miner `json:"miner"`
	TotalIncome           string `json:"totalIncome"`
	TotalExpenditure      string `json:"totalExpenditure"`
	BlockRewards          string `json:"blockRewards"`
	StorageDealPayments   string `json:"storageDealPayments"`
	RetrievalDealPayments string `json:"retrievalDealPayments"`
	NetworkFee            string `json:"networkFee"`
	Penalty               string `json:"penalty"`
	PreCommitDeposits     string `json:"preCommitDeposits"`
	InitialPledge         string `json:"initialPledge"`
	LockedFunds           string `json:"lockedFunds"`
	AvailableFunds        string `json:"availableFunds"`
}

type Fun struct {
	ID   string `json:"id"`
	Meta string `json:"meta"`
	Text string `json:"text"`
}

type Funds struct {
	PreCommitDeposits string `json:"preCommitDeposits"`
	InitialPledge     string `json:"initialPledge"`
	LockedFunds       string `json:"lockedFunds"`
	AvailableFunds    string `json:"availableFunds"`
}

type Income struct {
	Total                 string `json:"total"`
	BlockRewards          string `json:"blockRewards"`
	StorageDealPayments   string `json:"storageDealPayments"`
	RetrievalDealPayments string `json:"retrievalDealPayments"`
}

type Miner struct {
	ID                   string               `json:"id"`
	Address              string               `json:"address"`
	PeerID               string               `json:"peerId"`
	Owner                *Owner               `json:"owner"`
	Worker               *Worker              `json:"worker"`
	Name                 string               `json:"name"`
	Bio                  string               `json:"bio"`
	Location             string               `json:"location"`
	Contact              *Contact             `json:"contact"`
	Verified             bool                 `json:"verified"`
	ServiceDetails       *ServiceDetails      `json:"serviceDetails"`
	QualityIndicators    *QualityIndicators   `json:"qualityIndicators"`
	FinanceMetrics       *FinanceMetrics      `json:"financeMetrics"`
	AllServiceDetails    []*ServiceDetails    `json:"allServiceDetails"`
	AllQualityIndicators []*QualityIndicators `json:"allQualityIndicators"`
	AllFinanceMetrics    []*FinanceMetrics    `json:"allFinanceMetrics"`
	StorageDeal          *StorageDeal         `json:"storageDeal"`
	Transaction          *Transaction         `json:"transaction"`
	Sector               *Sector              `json:"sector"`
	Penalty              *Penalty             `json:"penalty"`
	Deadline             *Deadline            `json:"deadline"`
	StorageDeals         []*StorageDeal       `json:"storageDeals"`
	Transactions         []*Transaction       `json:"transactions"`
	Sectors              []*Sector            `json:"sectors"`
	Penalties            []*Penalty           `json:"penalties"`
	Deadlines            []*Deadline          `json:"deadlines"`
}

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type Owner struct {
	ID                  string   `json:"id"`
	Miners              []*Miner `json:"miners"`
	Address             string   `json:"address"`
	Actor               Actor    `json:"actor"`
	Balance             string   `json:"balance"`
	Messages            int      `json:"messages"`
	CreatedAt           int64    `json:"createdAt"`
	LatestTransactionAt int64    `json:"latestTransactionAt"`
}

type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
}

type Penalty struct {
	ID        string       `json:"id"`
	Fee       string       `json:"fee"`
	Type      *PenaltyType `json:"type"`
	Height    int64        `json:"height"`
	Timestamp time.Time    `json:"timestamp"`
}

type QualityIndicators struct {
	// ID                   string `json:"id"`
	Miner                *Miner `json:"miner"`
	QualityAdjPower      string `json:"qualityAdjPower"`
	RawBytePower         string `json:"rawBytePower"`
	QualityAdjPowerRatio string `json:"qualityAdjPowerRatio"`
	RawBytePowerRatio    string `json:"rawBytePowerRatio"`
	WinCount             uint64 `json:"winCount"`
	FaultySectors        uint64 `json:"faultySectors"`
	DataStored           string `json:"dataStored"`
	BlocksMined          uint64 `json:"blocksMined"`
	FeeDebt              string `json:"feeDebt"`
	MiningEfficiency     uint64 `json:"miningEfficiency"`
}

type Sector struct {
	ID              string       `json:"id"`
	Miner           *Miner       `json:"miner"`
	Size            string       `json:"size"`
	ActivationEpoch int64        `json:"activationEpoch"`
	ExpirationEpoch int64        `json:"expirationEpoch"`
	State           *SectorState `json:"state"`
	InitialPledge   string       `json:"initialPledge"`
	Faults          []*Fault     `json:"faults"`
}

type ServiceDetails struct {
	// ID                string `json:"id"`
	Miner             *Miner `json:"miner"`
	Storage           bool   `json:"storage"`
	Retrieval         bool   `json:"retrieval"`
	Repair            bool   `json:"repair"`
	OnlineDeals       bool   `json:"onlineDeals"`
	OfflineDeals      bool   `json:"offlineDeals"`
	StorageAskPrice   string `json:"storageAskPrice"`
	RetrievalAskPrice string `json:"retrievalAskPrice"`
	MinPieceSize      uint64 `json:"minPieceSize"`
	MaxPieceSize      uint64 `json:"maxPieceSize"`
}

type StorageDeal struct {
	ID                int    `json:"id"`
	Miner             *Miner `json:"miner"`
	MessageID         string `json:"messageId"`
	ClientID          string `json:"clientId"`
	ProviderID        string `json:"providerId"`
	ClientAddress     string `json:"clientAddress"`
	Price             string `json:"price"`
	StartEpoch        int64  `json:"startEpoch"`
	EndEpoch          int64  `json:"endEpoch"`
	Duration          int    `json:"duration"`
	PaddedPieceSize   uint64 `json:"paddedPieceSize"`
	UnPaddedPieceSize uint64 `json:"unPaddedPieceSize"`
	PieceCid          string `json:"pieceCID"`
	Verified          bool   `json:"verified"`
}

type Transaction struct {
	ID              string    `json:"id"`
	Miner           *Miner    `json:"miner"`
	TransactionType string    `json:"transactionType"`
	Amount          string    `json:"amount"`
	Sender          string    `json:"sender"`
	Receiver        string    `json:"receiver"`
	Height          int64     `json:"height"`
	Timestamp       time.Time `json:"timestamp"`
	NetworkFee      string    `json:"networkFee"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Worker struct {
	ID                  string `json:"id"`
	Miner               *Miner `json:"miner"`
	Address             string `json:"address"`
	Actor               Actor  `json:"actor"`
	Balance             string `json:"balance"`
	Messages            int    `json:"messages"`
	CreatedAt           int64  `json:"createdAt"`
	LatestTransactionAt int64  `json:"latestTransactionAt"`
}

type Actor string

const (
	ActorStorageMiner   Actor = "STORAGE_MINER"
	ActorRetrievalMiner Actor = "RETRIEVAL_MINER"
	ActorAccount        Actor = "ACCOUNT"
)

var AllActor = []Actor{
	ActorStorageMiner,
	ActorRetrievalMiner,
	ActorAccount,
}

func (e Actor) IsValid() bool {
	switch e {
	case ActorStorageMiner, ActorRetrievalMiner, ActorAccount:
		return true
	}
	return false
}

func (e Actor) String() string {
	return string(e)
}

func (e *Actor) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Actor(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Actor", str)
	}
	return nil
}

func (e Actor) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FaultType string

const (
	FaultTypeDeclared FaultType = "DECLARED"
	FaultTypeDetected FaultType = "DETECTED"
	FaultTypeSkipped  FaultType = "SKIPPED"
)

var AllFaultType = []FaultType{
	FaultTypeDeclared,
	FaultTypeDetected,
	FaultTypeSkipped,
}

func (e FaultType) IsValid() bool {
	switch e {
	case FaultTypeDeclared, FaultTypeDetected, FaultTypeSkipped:
		return true
	}
	return false
}

func (e FaultType) String() string {
	return string(e)
}

func (e *FaultType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FaultType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FaultType", str)
	}
	return nil
}

func (e FaultType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PenaltyType string

const (
	PenaltyTypePrecommitExpiryPenalty  PenaltyType = "PRECOMMIT_EXPIRY_PENALTY"
	PenaltyTypeUndeclaredFaultPenalty  PenaltyType = "UNDECLARED_FAULT_PENALTY"
	PenaltyTypeDeclaredFaultPenalty    PenaltyType = "DECLARED_FAULT_PENALTY"
	PenaltyTypeOngoingFaultPenalty     PenaltyType = "ONGOING_FAULT_PENALTY"
	PenaltyTypeTerminationFaultPenalty PenaltyType = "TERMINATION_FAULT_PENALTY"
	PenaltyTypeConsensusFaultPenalty   PenaltyType = "CONSENSUS_FAULT_PENALTY"
)

var AllPenaltyType = []PenaltyType{
	PenaltyTypePrecommitExpiryPenalty,
	PenaltyTypeUndeclaredFaultPenalty,
	PenaltyTypeDeclaredFaultPenalty,
	PenaltyTypeOngoingFaultPenalty,
	PenaltyTypeTerminationFaultPenalty,
	PenaltyTypeConsensusFaultPenalty,
}

func (e PenaltyType) IsValid() bool {
	switch e {
	case PenaltyTypePrecommitExpiryPenalty, PenaltyTypeUndeclaredFaultPenalty, PenaltyTypeDeclaredFaultPenalty, PenaltyTypeOngoingFaultPenalty, PenaltyTypeTerminationFaultPenalty, PenaltyTypeConsensusFaultPenalty:
		return true
	}
	return false
}

func (e PenaltyType) String() string {
	return string(e)
}

func (e *PenaltyType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PenaltyType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PenaltyType", str)
	}
	return nil
}

func (e PenaltyType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SectorState string

const (
	SectorStatePrecommited SectorState = "PRECOMMITED"
	SectorStateCommitted   SectorState = "COMMITTED"
	SectorStateActive      SectorState = "ACTIVE"
	SectorStateFaulty      SectorState = "FAULTY"
	SectorStateRecovering  SectorState = "RECOVERING"
	SectorStateTerminated  SectorState = "TERMINATED"
)

var AllSectorState = []SectorState{
	SectorStatePrecommited,
	SectorStateCommitted,
	SectorStateActive,
	SectorStateFaulty,
	SectorStateRecovering,
	SectorStateTerminated,
}

func (e SectorState) IsValid() bool {
	switch e {
	case SectorStatePrecommited, SectorStateCommitted, SectorStateActive, SectorStateFaulty, SectorStateRecovering, SectorStateTerminated:
		return true
	}
	return false
}

func (e SectorState) String() string {
	return string(e)
}

func (e *SectorState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SectorState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SectorState", str)
	}
	return nil
}

func (e SectorState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Sort string

const (
	SortAsc  Sort = "asc"
	SortDesc Sort = "desc"
)

var AllSort = []Sort{
	SortAsc,
	SortDesc,
}

func (e Sort) IsValid() bool {
	switch e {
	case SortAsc, SortDesc:
		return true
	}
	return false
}

func (e Sort) String() string {
	return string(e)
}

func (e *Sort) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Sort(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Sort", str)
	}
	return nil
}

func (e Sort) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TransactionDirection string

const (
	TransactionDirectionIncoming TransactionDirection = "INCOMING"
	TransactionDirectionOutgoing TransactionDirection = "OUTGOING"
)

var AllTransactionDirection = []TransactionDirection{
	TransactionDirectionIncoming,
	TransactionDirectionOutgoing,
}

func (e TransactionDirection) IsValid() bool {
	switch e {
	case TransactionDirectionIncoming, TransactionDirectionOutgoing:
		return true
	}
	return false
}

func (e TransactionDirection) String() string {
	return string(e)
}

func (e *TransactionDirection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TransactionDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TransactionDirection", str)
	}
	return nil
}

func (e TransactionDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TransactionType string

const (
	TransactionTypeStorageDeal   TransactionType = "STORAGE_DEAL"
	TransactionTypeRetrievalDeal TransactionType = "RETRIEVAL_DEAL"
	TransactionTypeBlockReward   TransactionType = "BLOCK_REWARD"
	TransactionTypeNetworkFee    TransactionType = "NETWORK_FEE"
	TransactionTypePenalty       TransactionType = "PENALTY"
)

var AllTransactionType = []TransactionType{
	TransactionTypeStorageDeal,
	TransactionTypeRetrievalDeal,
	TransactionTypeBlockReward,
	TransactionTypeNetworkFee,
	TransactionTypePenalty,
}

func (e TransactionType) IsValid() bool {
	switch e {
	case TransactionTypeStorageDeal, TransactionTypeRetrievalDeal, TransactionTypeBlockReward, TransactionTypeNetworkFee, TransactionTypePenalty:
		return true
	}
	return false
}

func (e TransactionType) String() string {
	return string(e)
}

func (e *TransactionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TransactionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TransactionType", str)
	}
	return nil
}

func (e TransactionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
