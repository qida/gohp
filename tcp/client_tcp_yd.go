package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	RequestType开始 = "startAudioStream"
	RequestType停止 = "stopAudioStream"
)

type Hand struct {
	Version     int    `json:"version"`
	RequestType string `json:"requestType"`
	RequestID   int64  `json:"requestId"`
	App         string `json:"app"`
	Stream      string `json:"stream"`
}

func newTCPCMD(request_type string, stream_id string) Hand {
	return Hand{
		Version:     1,
		RequestType: request_type,
		RequestID:   time.Now().Unix(),
		App:         "rtp",
		Stream:      stream_id,
	}
}

type TCP struct {
	ReceiveChan chan []byte
	Conn        net.Conn

	streamId string
}

func (t *TCP) SendStart(stream_id string) (_err error) {
	cmd := newTCPCMD(RequestType开始, stream_id)
	t.streamId = stream_id
	bytCmd, _ := json.Marshal(cmd)
	_, _err = fmt.Fprint(t.Conn, string(bytCmd))
	return
}

func (t *TCP) SendStop() (_err error) {
	if t.Conn == nil {
		_err = fmt.Errorf("链接无效")
		return
	}
	if t.streamId == "" {
		_err = fmt.Errorf("流无效")
		return
	}
	cmd := newTCPCMD(RequestType停止, t.streamId)
	bytCmd, _ := json.Marshal(cmd)
	_, _err = fmt.Fprint(t.Conn, string(bytCmd))
	return
}

func (t *TCP) Close() {
	if t.Conn != nil {
		t.Conn.Close()
		close(t.ReceiveChan)
	}
}
func NewTcpClient() (tcp *TCP) {
	tcp = new(TCP)
	tcp.ReceiveChan = make(chan []byte)
	return
}

func (t *TCP) Link(tcp_link string) (err error) {
	t.Conn, err = net.Dial("tcp", tcp_link)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	go func() {
		var buf = make([]byte, 1024*2)
		for {
			buf = make([]byte, 1024*2)
			n, err := t.Conn.Read(buf)
			if err != nil {
				fmt.Printf("客户端已退出..\n")
				return
			}
			t.ReceiveChan <- buf[:n]
		}
	}()
	return
}
