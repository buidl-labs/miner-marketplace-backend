package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
