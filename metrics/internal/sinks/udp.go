package sinks

import (
	"log"
	"net"
)

type UdpClient struct {
	Endpoint Endpoint
}

func NewUdpClient(endpoint Endpoint) *UdpClient {
	return &UdpClient{Endpoint: endpoint}
}

func (u *UdpClient) SendMessage(message []byte) error {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(u.Endpoint.Host, string(u.Endpoint.Port)))
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v", err)
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Failed to dial UDP: %v", err)
		return err
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		log.Printf("Failed to send UDP message: %v", err)
		return err
	}

	log.Println("Message sent via UDP.")
	return nil
}
