package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetJson(url string, target interface{}) (interface{}, error) {
	r, err := http.Get(url)
	if err != nil {
		fmt.Println("err:", err)
		return target, err
	}
	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		fmt.Println("readErr:", readErr)
		return target, readErr
	}

	jsonErr := json.Unmarshal(body, target)
	if jsonErr != nil {
		fmt.Println("jsonErr:", jsonErr)
		return target, jsonErr
	}
	return target, nil
}

func GenerateTransactionTypesQuery(transactionTypes []bool) string {
	fmt.Println("len", len(transactionTypes))
	query := ""
	for i, tt := range transactionTypes {
		switch i {
		case 0:
			if tt {
				query += "transaction_type = 'Collateral Deposit' OR "
			}
		case 1:
			if tt {
				query += "transaction_type = 'Block Reward' OR "
			}
		case 2:
			if tt {
				query += "transaction_type = 'Deals Publish' OR "
			}
		case 3:
			if tt {
				query += "transaction_type = 'Penalty' OR "
			}
		case 4:
			if tt {
				query += "transaction_type = 'Transfer' OR "
			}
		case 5:
			if tt {
				query += "transaction_type = 'Others' OR "
			}
		}
	}
	fmt.Println("tempquery:", query)
	if len(query) > 0 && strings.HasSuffix(query, " OR ") {
		query = query[:len(query)-len(" OR ")]
	}
	fmt.Println("finalquery:", query)
	return query
}
