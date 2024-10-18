package sinks

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type TcpClient struct {
	Endpoint Endpoint
	Conn     net.Conn
	mutex    sync.Mutex
}

const (
	maxRetries = 10
	retryDelay = 5 * time.Second
)

func NewTcpClient(endpoint Endpoint) *TcpClient {
	return &TcpClient{Endpoint: endpoint}
}

func (c *TcpClient) InitialConnect() error {

	var err error
	addr := fmt.Sprintf("%s:%s", c.Endpoint.Host, c.Endpoint.Port)
	for i := 0; i < maxRetries; i++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			c.Conn = conn
			log.Printf("TcpClient connected to %s", addr)
			return nil
		}
		fmt.Printf("Waiting for CloudWatch Agent to be reachable... (%d/%d)\n", i+1, maxRetries)
		time.Sleep(retryDelay)
	}

	log.Printf("Failed to connect: %v", err)
	return err
}

func (c *TcpClient) Warmup() error {
	return c.establishConnection()
}

func (c *TcpClient) SendMessage(message []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.waitForOpenConnection(); err != nil {
		return err
	}

	_, err := c.Conn.Write(message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	return nil
}

func (c *TcpClient) Disconnect(reason string) {
	log.Printf("TcpClient disconnected due to: %s", reason)
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (c *TcpClient) waitForOpenConnection() error {
	if c.Conn == nil {
		return c.establishConnection()
	}
	return nil
}

func (c *TcpClient) establishConnection() error {
	if c.Conn != nil {
		return nil
	}
	return c.InitialConnect()
}
