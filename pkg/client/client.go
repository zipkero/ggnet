package client

import (
	"encoding/binary"
	"errors"
	"github.com/zipkero/ggnet/pkg/message"
	"log"
	"net"
	"sync"
)

type Client struct {
	host      string
	port      string
	conn      net.Conn
	SendCh    chan message.Message
	ReceiveCh chan message.Message
}

func NewClient(host, port string) *Client {
	return &Client{
		host:      host,
		port:      port,
		SendCh:    make(chan message.Message),
		ReceiveCh: make(chan message.Message),
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
	return nil
}

func (c *Client) Listen() error {
	var wg sync.WaitGroup
	wg.Add(2)

	go c.receive(&wg)
	go c.send(&wg)

	wg.Wait()

	return nil
}

func (c *Client) receive(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		lengthBuffer := make([]byte, 4)
		_, err := c.conn.Read(lengthBuffer)
		if err != nil {
			log.Println(err)
		}

		messageLength := binary.BigEndian.Uint32(lengthBuffer)
		messageBuffer := make([]byte, messageLength)

		_, err = c.conn.Read(messageBuffer)
		if err != nil {
			log.Println(err)
		}

		messageType := binary.BigEndian.Uint16(messageBuffer[:2])
		messageContent := string(messageBuffer[2:])

		c.ReceiveCh <- message.Message{
			Type:    messageType,
			Content: messageContent,
		}
	}
}

func (c *Client) send(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case msg := <-c.SendCh:
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
