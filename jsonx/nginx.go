package jsonx

import (
	"log"

	"github.com/recoye/config"
)

type Upstream struct {
	Server string
}

type ServConf struct {
	ServerName string
}

func ReadNginxConf(file_path string) (_err error) {
	conf := config.New(file_path)
	co := &ServConf{}
	err := conf.Unmarshal(co)
	if err == nil {
		log.Println(co)
	} else {
		log.Println(err)
	}
	return
}
