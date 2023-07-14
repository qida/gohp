package logx

import (
	"fmt"
	"time"

	"gohp/tcp"

	"github.com/axgle/mahonia"
	"github.com/beego/beego/v2/adapter/logs"
	beego "github.com/beego/beego/v2/server/web"
)

var (
	DebugList map[string]*tcp.Client
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

func ServerDebug(port int) {
	fmt.Printf("调试 在 %d 监听...\r\n", port)
	DebugList = make(map[string]*tcp.Client)
	server := tcp.New(fmt.Sprintf("0.0.0.0:%d", port), "")
	// utf-8=>gb18030
	//dec := mahonia.NewDecoder("GB18030")
	// gb18030=>utf-8
	//enc := mahonia.Newutil.Encoder("GB18030")
	server.OnNewClient(func(c *tcp.Client) {
		fmt.Printf("新的调试端\r\n")
		c.Send(fmt.Sprintf("Welcome %s \n", c.GetConn().RemoteAddr().String()))
	})
	server.OnNewMessage(func(c *tcp.Client, message string) {
		if message == "debug\r\n" {
			DebugList[c.GetConn().RemoteAddr().String()] = c
			c.Send("Welcome Debugger\r\n")
		} else if message == "\n" {
			//不处理
		} else {
			for _, v := range DebugList {
				v.Send(message)
			}
		}
	})
	server.OnClientConnectionClosed(func(c *tcp.Client, err error) {
		fmt.Printf("调试端断开\r\n")
		delete(DebugList, c.GetConn().RemoteAddr().String())
	})
	server.Listen()
}
