package httpx

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ClientMangerCenter struct {
	Register    chan *ClientWS
	Unregister  chan *ClientWS
	Broadcaster chan *Broadcaster

	MapWSNet
	MapWSSub
}

type ClientWS struct {
	KeyNet string
	KeySub string

	Conn *websocket.Conn
	Link *websocket.Conn
	Url  *url.URL

	MsgConn   chan []byte `json:"-"`
	MsgLink   chan []byte `json:"-"`
	SubDevice SubDevice
}

type SubDevice struct {
	DeviceId string
	IdScheme int

	ZoneId string
}

type Broadcaster struct {
	KeySub  string
	Message []byte
}

var clientMangerCenter = ClientMangerCenter{
	Register:    make(chan *ClientWS),
	Unregister:  make(chan *ClientWS),
	Broadcaster: make(chan *Broadcaster),

	MapWSNet: MapWSNet{MapClientWSNets: make(map[string]*ClientWS), Lock: &sync.RWMutex{}},
	MapWSSub: MapWSSub{MapClientWSSubs: make(map[string]map[string]*ClientWS), Lock: &sync.RWMutex{}},
}

func NewWsClientMangerCenter() *ClientMangerCenter {
	return &clientMangerCenter
}

func (t *SubDevice) GetKeySub() string {
	return fmt.Sprintf("%s_%d\r\n", t.DeviceId, t.IdScheme)
}

func NewClientWS(conn *websocket.Conn, sub_device SubDevice) *ClientWS {
	client := &ClientWS{
		KeyNet: conn.RemoteAddr().String(),
		KeySub: sub_device.GetKeySub(),

		Conn:      conn,
		MsgConn:   make(chan []byte),
		SubDevice: sub_device,
	}
	go client.Read()
	go client.SendMsgConn()
	go client.SendMsgLink()
	return client
}

func (t *ClientMangerCenter) Start() {
	for {
		select {
		case clientWS := <-t.Register:
			log.Printf("Register 收到:%s\r\n", clientWS.KeyNet)
			t.MapWSNet.Add(clientWS)
			t.MapWSSub.Add(clientWS)
			err := clientWS.Conn.WriteMessage(websocket.TextMessage, []byte("OK"))
			if err != nil {
				log.Printf("Register WriteMessage 错误:%+v\r\n", err)
				clientWS.Close()
				continue
			}
			log.Printf("Register KeyNet:%s 全部终端: %d\r\n", clientWS.KeyNet, len(t.MapClientWSNets))
			log.Printf("Register KeySub:%s 指定订阅: %d\r\n", clientWS.KeySub, len(t.MapClientWSSubs[clientWS.KeySub]))
			if clientWS.Url == nil {
				continue
			}
			go clientWS.StartLink(t)
		case clientWS := <-t.Unregister:
			log.Printf("Unregister 收到:%s\r\n", clientWS.KeyNet)
			clientWS.Conn.Close()
			if clientWS.Link != nil {
				clientWS.Link.Close()
			}
			t.MapWSSub.Delete(clientWS)
			t.MapWSNet.Delete(clientWS)
			log.Printf("Unregister KeyNet:%s 全部终端: %d\r\n", clientWS.KeyNet, len(t.MapClientWSNets))
			log.Printf("Unregister KeySub:%s 指定订阅: %d\r\n", clientWS.KeySub, len(t.MapClientWSSubs[clientWS.KeySub]))
		case message := <-t.Broadcaster:
			t.MapWSSub.Lock.RLock()
			if clientWSs, ok := t.MapClientWSSubs[message.KeySub]; ok {
				for _, clientWS := range clientWSs {
					select {
					case <-time.After(time.Second * 3):
						log.Printf("Broadcaster 超时 客户端:%s\r\n", clientWS.KeyNet)
						clientWS.Close()
					case clientWS.MsgConn <- message.Message:
						log.Printf("Broadcaster KeyNet:%s\r\n", clientWS.KeyNet)
					default:
						log.Printf("Broadcaster KeyNet:%s 全部终端: %d\r\n", clientWS.KeyNet, len(t.MapClientWSNets))
						log.Printf("Broadcaster KeySub:%s 指定订阅: %d\r\n", clientWS.KeySub, len(t.MapClientWSSubs[clientWS.KeySub]))
					}
				}
				log.Printf("WS广播 Tag 订阅:%s 数量:%d\r\n", message.KeySub, len(clientWSs))
			}
			t.MapWSSub.Lock.RUnlock()
		}
	}
}

func (t *ClientMangerCenter) Broadcast(sub_device SubDevice, message []byte) {
	broadcaster := Broadcaster{
		KeySub:  sub_device.GetKeySub(),
		Message: message,
	}
	t.Broadcaster <- &broadcaster

	// 特殊全部订阅
	broadcasterAll := Broadcaster{
		KeySub:  "all_-1\r\n",
		Message: message,
	}
	t.Broadcaster <- &broadcasterAll
}

func (c *ClientWS) Read() {
	defer func() {
		log.Println("ClientWS Read Over")
		c.Close()
	}()
	log.Printf("开始 Read : %s\r\n", c.KeyNet)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("IsUnexpectedCloseError 错误: %+v \r\n", err)
				return
			}
			log.Printf("ReadMessage 错误: %+v \r\n", err)
			return
		}
		log.Printf("Read message: %s\r\n", message)
	}
}

func (c *ClientWS) SendMsgConn() (_err error) {
	defer func() {
		log.Println("ClientWS Send Over")
		c.Close()
	}()
	log.Printf("开始 Send : %s\r\n", c.KeyNet)
	for message := range c.MsgConn {
		if _err = c.Conn.WriteMessage(websocket.TextMessage, message); _err != nil {
			log.Printf("WriteMessage 错误\r\n", zap.Error(_err))
			return
		}
		log.Printf("发送Send message to %s\r\n", c.KeyNet)
	}
	return
}

func (c *ClientWS) SendMsgLink() (_err error) {
	defer func() {
		log.Println("ClientWS Send Over")
		c.Close()
	}()
	log.Printf("开始 Send : %s\r\n", c.KeyNet)
	for message := range c.MsgLink {
		if _err = c.Link.WriteMessage(websocket.TextMessage, message); _err != nil {
			log.Printf("WriteMessage 错误\r\n", zap.Error(_err))
			return
		}
		log.Printf("Send message to %s\r\n", c.KeyNet)
	}
	return
}

func (c *ClientWS) StartLink(t *ClientMangerCenter) {
	defer func() {
		log.Println("ClientLink  Over")
		c.Close()
	}()
	log.Printf("建立与区域连接 %s\r\n", c.Url.String())
	var conn *websocket.Conn
	var err error
	for {
		conn, _, err = websocket.DefaultDialer.Dial(c.Url.String(), nil)
		if err != nil {
			log.Printf("连接WebSocket服务器失败 %+v\r\n", err)
			time.Sleep(time.Second * 20)
			continue
		}
		break
	}
	defer conn.Close()
	c.Link = conn
	log.Printf("读取 连接区域 [%s] WebSocket服务器成功\r\n", c.Url.String())
	for {
		_, message, err := c.Link.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("IsUnexpectedCloseError 错误:%+v\r\n", err)
				return
			}
			log.Printf("ReadMessage 错误:%+v\r\n", err)
			return
		}
		log.Printf("Read message: %s\r\n", message)
		t.Broadcast(c.SubDevice, message)
	}
}

func (t *ClientMangerCenter) ListSub() []string {
	var list []string
	for key := range t.MapClientWSSubs {
		list = append(list, key)
	}
	return list
}

func (t *ClientMangerCenter) CountSub() int {
	return len(t.MapClientWSSubs)
}

func (t *ClientMangerCenter) CountNet() int {
	return len(t.MapClientWSNets)
}

func (t *ClientMangerCenter) ListNet() []string {
	var list []string
	for key := range t.MapClientWSNets {
		list = append(list, key)
	}
	return list
}

func (c *ClientWS) Close() {
	log.Printf("ClientWS Close : %s\r\n", c.KeyNet)
	clientMangerCenter.Unregister <- c
}
