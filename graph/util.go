package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buidl-labs/filecoin-chain-indexer/model/messages"
)

type OwnerChange struct {
	Epoch int64  `json:"epoch,omitempty"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
}

type WorkerChange struct {
	Epoch int64  `json:"epoch,omitempty"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
}

type ControlChange struct {
	Epoch int64    `json:"epoch,omitempty"`
	From  []string `json:"from,omitempty"`
	To    []string `json:"to,omitempty"`
}

type MinerAddressChanges struct {
	OwnerChanges   []OwnerChange   `json:"ownerChanges,omitempty"`
	WorkerChanges  []WorkerChange  `json:"workerChanges,omitempty"`
	ControlChanges []ControlChange `json:"controlChanges,omitempty"`
}

func GetMinerAddressChanges() map[string]MinerAddressChanges {
	url := os.Getenv("ADDR_CHANGES_URL")
	myClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("URLLL", url, " ERRORRR", err)
		panic(err)
	}
	req.Header.Set("User-Agent", "miner-marketplace-app")
	res, getErr := myClient.Do(req)
	if getErr != nil {
		fmt.Println("URLLL", url, " ERRORRR", err)
		panic(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		panic(readErr)
	}
	fmt.Println("body", body)

	var minerAddressChanges map[string]MinerAddressChanges
	jsonErr := json.Unmarshal(body, &minerAddressChanges)
	if jsonErr != nil {
		panic(jsonErr)
	}
	fmt.Println(minerAddressChanges)
	return minerAddressChanges
}

func ComputeBigIntSum(total *big.Int, amounts []string) *big.Int {
	for _, amount := range amounts {
		n := new(big.Int)
		n, ok := n.SetString(amount, 10)
		if !ok {
			fmt.Println("SetString: error")
		}
		fmt.Println(n)
		total.Add(total, n)
	}
	return total
}

func GetTransactionType(methodName string) string {
	switch methodName {
	case "Send":
		return "DEAL"
	case "ApplyRewards":
		return "REWARD"
	default:
		return "NETWORK_FEE"
	}
	// return ""
}

func DeriveTransactionLabels(transaction messages.Transaction) (string, string, bool) {
	label := transaction.MethodName
	direction := "+"
	gas := false
	methodName := transaction.MethodName
	actorName := transaction.ActorName
	if strings.HasPrefix(actorName, "storageMinerActor") {
		switch methodName {
		case "SubmitWindowedPoSt":
			direction = "-"
			gas = true
		case "PreCommitSector":
			label = "PreCommit Deposit"
			direction = "collateral(-)"
			gas = true
		case "ProveCommitSector":
			label = "Initial Pledge"
			direction = "collateral(-)"
			gas = true
		case "ApplyRewards":
			label = "Block Reward"
			direction = "+"
		}
	}
	return label, direction, gas
}

func CalculateDealPrice(pricePerEpoch string, startEpoch int64, endEpoch int64) (string, error) {
	pricePerEpochInt, err := strconv.ParseInt(pricePerEpoch, 10, 64)
	if err != nil {
		return "0", err
	}
	totalPrice := (endEpoch - startEpoch) * pricePerEpochInt
	return fmt.Sprintf("%v", totalPrice), nil
}

type AlphaBetaFilter struct {
	alpha big.Int // Q.128
	beta  big.Int // Q.128
	p     *big.Int
	v     *big.Int
}

func NewAlphaBetaFilter(p, v *big.Int) AlphaBetaFilter {
	return AlphaBetaFilter{
		p: p,
		v: v,
	}
}

func ProjectFutureReward(days int, sectorQAP, nwqapP, nwqapV, perEpochRewardP, perEpochRewardV *big.Int) *big.Int {
	networkQAPFilter := NewAlphaBetaFilter(nwqapP, nwqapV)
	perEpochRewardFilter := NewAlphaBetaFilter(perEpochRewardP, perEpochRewardV)

	return new(big.Int).Mul(sectorQAP, ExtrapolateCumsumRatio(perEpochRewardFilter, networkQAPFilter, days*2880))
}

func ExtrapolateCumsumRatio(numerator, denominator AlphaBetaFilter, futureT int) *big.Int {
	v1, _ := new(big.Float).SetInt(
		new(big.Int).Add(denominator.p, denominator.v),
	).Float64()
	x2a := math.Log(v1)
	v2, _ := new(big.Float).SetInt(
		new(big.Int).Add(
			new(big.Int).Add(denominator.p, denominator.v),
			new(big.Int).Mul(denominator.v, big.NewInt(int64(futureT))),
		),
	).Float64()
	x2b := math.Log(v2)
	m1 := new(big.Int).Mul(new(big.Int).Mul(denominator.v, numerator.p), big.NewInt(int64(x2b-x2a)))
	m2 := new(big.Int).Mul(
		numerator.v,
		new(big.Int).Add(
			new(big.Int).Mul(denominator.p, big.NewInt(int64(x2a-x2b))),
			new(big.Int).Mul(denominator.v, big.NewInt(int64(futureT))),
		),
	)
	return new(big.Int).Div(new(big.Int).Add(m1, m2), denominator.v.Exp(denominator.v, big.NewInt(2), nil))
}
