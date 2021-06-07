package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetJson(url string, target interface{}) (interface{}, error) {
	r, err := http.Get(url)
	if err != nil {
		return target, err
	}
	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return target, readErr
	}

	jsonErr := json.Unmarshal(body, target)
	if jsonErr != nil {
		return target, jsonErr
	}
	return target, nil
}
