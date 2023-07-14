package logx

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var config = `
caller: true
console: true
level: "debug"
logs:
  path: "logs/service-manager.log"
  maxAge: 30
  size: 5
  isGzip: true
file:
  - path: "logs/debug.log"
    level: "debug"
    size: 500
  - path: "logs/info.log"
    level: "info"
    maxAge: 30
`

func TestLogger(t *testing.T) {
	conf := make(map[string]interface{})
	if err := yaml.NewDecoder(strings.NewReader(config)).Decode(&conf); err != nil {
		t.Errorf("yaml decode error: %s", err)
		return
	}
	err := DefaultWithMap(conf)
	if err != nil {
		t.Errorf("new logger error: %s", err)
		return
	}
	logger.Debug("configs", zap.String("config", config))
	logger.Warn("this is a test!!")
	logger.Info("logger test start:")
	for i := 0; i < 10; i++ {
		logger.Info("time now", zap.Int("times", i), zap.Time("now", time.Now()))
		time.Sleep(1 * time.Second)
	}
	logger.Info("logger test done.")
}

func TestBigFileLogger(t *testing.T) {
	conf := make(map[string]interface{})
	if err := yaml.NewDecoder(strings.NewReader(config)).Decode(&conf); err != nil {
		t.Errorf("yaml decode error: %s", err)
		return
	}
	err := DefaultWithMap(conf)
	if err != nil {
		t.Errorf("new logger error: %s", err)
		return
	}
	logger.Debug("configs", zap.String("config", config))
	logger.Warn("this is a test!!")
	logger.Info("logger test start:")
	pi, err := os.ReadFile("example/pi.txt")
	if err != nil {
		logger.Error("ReadFile error", zap.Error(err))
		return
	}
	for i := 0; i < 1000; i++ {
		logger.Info("pi", zap.Int("times", i), zap.ByteString("Π", pi))
	}
	logger.Info("logger test done.")
}

func TestLogger1(t *testing.T) {

	//log.Println("这是一条优雅的日志。")
	//v := "优雅的"
	//log.Printf("这是一个%s日志\n", v)
	////fatal系列函数会在写入日志信息后调用os.Exit(1)。Panic系列函数会在写入日志信息后panic。
	//log.Fatalln("这是一天会触发fatal的日志")
	//log.Panicln("这是一个会触发panic的日志。") //执行后会自动触发一个异常
	//
	log.SetPrefix("[PS]")
	log.Println("这是一条很普通的日志。")

	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("这是一条优雅的日志。")

}
