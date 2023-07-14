package logx

import (
	"fmt"
	"time"

	"github.com/axgle/mahonia"
	"github.com/beego/beego/v2/adapter/logs"
	beego "github.com/beego/beego/v2/server/web"
)


var (
	LogConn = logs.NewLogger(1000)
	LogTcp  = logs.NewLogger(1000)
	LogMail = logs.NewLogger(1000)
	Enc     = mahonia.NewEncoder("gb18030")
)

func ClientDebug(ip string, port int) {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async(1e3)
	LogTcp.SetLevel(logs.LevelDebug)
	LogTcp.SetLogger(logs.AdapterConn, fmt.Sprintf(`{"net":"tcp","addr":"%s:%d","reconnect":true}`, ip, port))
	LogConn.SetLevel(logs.LevelDebug)
	LogConn.SetLogger(logs.AdapterConsole)
}

func Email() {
	LogMail.Async()
	LogMail.EnableFuncCallDepth(true)
	err := LogMail.SetLogger(logs.AdapterMail, `{"level":7,"username":"sunqida@126.com","password":"","fromAddress":"sunqida@126.com","subject":"", "host":"smtp.126.com:994","sendTos":["sunqida@foxmail.com"]}`) //654/994
	if err != nil {
		panic(err.Error())
	}
	if beego.BConfig.RunMode == "dev" {
		LogMail.Notice("Api Test系统开始运行：%v", time.Now())
	} else {
		LogMail.Notice("Api Prod系统开始运行：%v", time.Now())
	}
}
