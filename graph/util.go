package graph

import (
	"encoding/json"
	"fmt"
	"strconv"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"time"
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

func CalculateDealPrice(pricePerEpoch string, startEpoch int64, endEpoch int64) (string, error) {
	pricePerEpochInt, err := strconv.ParseInt(pricePerEpoch, 10, 64)
	if err != nil {
		return "0", err
	}
	totalPrice := (endEpoch - startEpoch) * pricePerEpochInt
	return fmt.Sprintf("%v", totalPrice), nil
}
