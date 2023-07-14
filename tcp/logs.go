package tcp

import (
	"fmt"
)
var (
	DebugList map[string]*Client
)

func ServerDebug(port int) {
	fmt.Printf("调试 在 %d 监听...\r\n", port)
	DebugList = make(map[string]*Client)
	server := New(fmt.Sprintf("0.0.0.0:%d", port), "")
	// utf-8=>gb18030
	//dec := mahonia.NewDecoder("GB18030")
	// gb18030=>utf-8
	//enc := mahonia.Newutil.Encoder("GB18030")
	server.OnNewClient(func(c *Client) {
		fmt.Printf("新的调试端\r\n")
		c.Send(fmt.Sprintf("Welcome %s \n", c.GetConn().RemoteAddr().String()))
	})
	server.OnNewMessage(func(c *Client, message string) {
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
	server.OnClientConnectionClosed(func(c *Client, err error) {
		fmt.Printf("调试端断开\r\n")
		delete(DebugList, c.GetConn().RemoteAddr().String())
	})
	server.Listen()
}
