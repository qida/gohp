package jsonx

import (
	"fmt"
	"log"

	"github.com/recoye/config"
	"github.com/smartystreets/goconvey/web/server/parser"
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

func ReadNginxConf1(file_path string) (_err error) {
	p := parser.NewParser(file_path)
	c := p.Parse()
	fmt.Println(gonginx.DumpConfig(c, gonginx.IndentedStyle))
	return
}
