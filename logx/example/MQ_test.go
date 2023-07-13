package main

import (
	"encoding/json"
	"gobase/logger"
	"io/ioutil"
	"testing"
)

func Test_mq(t *testing.T) {

	jsonStr, err := ioutil.ReadFile("./mq.json")
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(jsonStr, &m)
	if err != nil {
		panic(err)
	}

	err = logger.DefaultWithMap(m)
	if err != nil {
		panic(err)
	}
	logger.Debug("白居易")

}
