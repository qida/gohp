package tcp

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *server
}

// TCP server
type server struct {
	address                  string // Address to open connection: localhost:9999
	relay                    string
	RelayMessage             chan string
	config                   *tls.Config
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

// Read client data from channel
func (c *Client) listen() {
	c.Server.onNewClientCallback(c)
	reader := bufio.NewReader(c.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
		c.Server.RelayMessage <- message
	}
}

// Get Conn
func (c *Client) GetConn() net.Conn {
	return c.conn
}

// Send text message to client
func (c *Client) Send(message string) error {
	return c.SendBytes([]byte(message))
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	if err != nil {
		c.conn.Close()
		c.Server.onClientConnectionClosed(c, err)
	}
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

// Creates new tcp server instance
func New(address, relay string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address:      address,
		relay:        relay,
		RelayMessage: make(chan string, 10),
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

func NewWithTLS(address, relay, certFile, keyFile string) *server {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error loading certificate files. Unable to create TCP server with TLS functionality.\r\n", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server := New(address, relay)
	server.config = config
	return server
}

func (s *server) initRelay() {
	var relayTcp net.Conn
	var err error
	for {
		msg := <-s.RelayMessage
		if s.relay != "" {
			if relayTcp == nil {
				relayTcp, err = net.Dial("tcp", s.relay)
				if err != nil {
					relayTcp = nil
					log.Println("Error starting TCP relay.", err)
					continue
				} else {
					defer relayTcp.Close()
				}
			}
			_, err = relayTcp.Write([]byte(msg))
			if err != nil {
				relayTcp = nil
				log.Println("Error TCP relay.", err)
				continue
			}
		}
	}
}
