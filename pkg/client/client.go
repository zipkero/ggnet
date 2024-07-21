package client

import (
	"encoding/binary"
	"errors"
	"github.com/zipkero/ggnet/pkg/message"
	"log"
	"net"
)

type Client struct {
	host string
	port string
	conn net.Conn
	Ch   chan message.Message
}

func NewClient(host, port string) *Client {
	return &Client{
		host: host,
		port: port,
		Ch:   make(chan message.Message),
	}
}

func (c *Client) Connect() error {
	ip, err := c.resolveAddress()
	if err != nil {
		return err
	}

	c.conn, err = net.Dial("tcp", net.JoinHostPort(ip, c.port))
	if err != nil {
		return err
	}

	go c.receive()
	go c.send()

	return nil
}

func (c *Client) receive() {
	for {

	}
}

func (c *Client) send() {
	for {
		select {
		case msg := <-c.Ch:
			var typeBytes = make([]byte, 2)
			binary.BigEndian.PutUint16(typeBytes, msg.Type)
			sendMessageBytes := append(typeBytes, []byte(msg.Content)...)

			lengthBuffer := len(sendMessageBytes)
			_, err := c.conn.Write([]byte{
				byte(lengthBuffer >> 24),
				byte(lengthBuffer >> 16),
				byte(lengthBuffer >> 8),
				byte(lengthBuffer),
			})

			if err != nil {
				log.Println(err)
			}

			_, err = c.conn.Write(sendMessageBytes)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (c *Client) resolveAddress() (string, error) {
	if ip := net.ParseIP(c.host); ip != nil {
		return ip.String(), nil
	}

	ips, err := net.LookupIP(c.host)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", errors.New("no ip address found")
	}

	return ips[0].String(), nil
}
