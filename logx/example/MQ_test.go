package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/qida/gohp/logx"
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

	err = logx.DefaultWithMap(m)
	if err != nil {
		panic(err)
	}
	logx.Debug("白居易")

}
