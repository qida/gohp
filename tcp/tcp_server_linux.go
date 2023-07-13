package tcp

import (
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/felixge/tcpkeepalive"
)

// // Start network server
// func (s *server) Listen() {
// 	go s.listenChannels()

// 	listener, err := net.Listen("tcp", s.address)
// 	if err != nil {
// 		log.Fatal("Error starting TCP server")
// 	}
// 	defer listener.Close()

// 	for {
// 		conn, _ := listener.Accept()
// 		kaConn, _ := tcpkeepalive.EnableKeepAlive(conn)
// 		kaConn.SetKeepAliveIdle(30 * time.Second)
// 		kaConn.SetKeepAliveCount(4)
// 		kaConn.SetKeepAliveInterval(5 * time.Second)
// 		s.joins <- kaConn
// 	}
// }

// Listen starts network server
func (s *server) Listen() {
	var listener net.Listener
	var err error
	if s.config == nil {
		listener, err = net.Listen("tcp", s.address)
	} else {
		listener, err = tls.Listen("tcp", s.address, s.config)
	}
	if err != nil {
		log.Fatal("Error starting TCP server.\r\n", err)
	}
	defer listener.Close()
	go s.initRelay()
	for {
		// conn, _ := listener.Accept()
		conn, _ := listener.Accept()
		kaConn, _ := tcpkeepalive.EnableKeepAlive(conn)
		kaConn.SetKeepAliveIdle(30 * time.Second)
		kaConn.SetKeepAliveCount(4)
		kaConn.SetKeepAliveInterval(5 * time.Second)
		client := &Client{
			conn:   kaConn,
			Server: s,
		}
		go client.listen()
	}
}
