package tcp

import (
	"net"
	"testing"
	"time"
)

func buildTestServer() *server {
	return New("localhost:9999", "")
}

func Test_accepting_new_client_callback(t *testing.T) {
	server := buildTestServer()

	newClient := make(chan struct{}, 1)
	messageReceived := make(chan string, 1)
	connectionClosed := make(chan struct{}, 1)

	server.OnNewClient(func(c *Client) {
		newClient <- struct{}{}
	})
	server.OnNewMessage(func(c *Client, message string) {
		messageReceived <- message
	})
	server.OnClientConnectionClosed(func(c *Client, err error) {
		connectionClosed <- struct{}{}
	})

	// Listen 会阻塞，需在 goroutine 中启动
	go server.Listen()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		t.Fatal("Failed to connect to test server")
	}
	if _, err := conn.Write([]byte("Test message\n")); err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	conn.Close()

	const timeout = 2 * time.Second

	select {
	case <-newClient:
	case <-time.After(timeout):
		t.Fatal("超时：未收到新客户端回调")
	}

	select {
	case msg := <-messageReceived:
		if msg != "Test message\n" {
			t.Errorf("收到消息 = %q, want %q", msg, "Test message\n")
		}
	case <-time.After(timeout):
		t.Fatal("超时：未收到消息回调")
	}

	select {
	case <-connectionClosed:
	case <-time.After(timeout):
		t.Fatal("超时：未收到连接关闭回调")
	}
}
