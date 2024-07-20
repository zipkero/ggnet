package client

import (
	"errors"
	"net"
)

type Client struct {
	host string
	port string
	conn net.Conn
}

func NewClient(host, port string) *Client {
	return &Client{
		host: host,
		port: port,
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
