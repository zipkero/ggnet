package main

import (
	"encoding/binary"
	"fmt"
	"github.com/zipkero/ggnet/internal/message"
	"log"
	"net"
)

func main() {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 5000,
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	sendMessage := &message.Message{
		Type: 1,
	}

	for {
		_, err = fmt.Scanln(&sendMessage.Content)
		if err != nil {
			log.Println(err)
		}
		var typeBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(typeBytes, sendMessage.Type)
		sendMessageBytes := append(typeBytes, []byte(sendMessage.Content)...)

		lengthBuffer := len(sendMessageBytes)
		_, err = conn.Write([]byte{
			byte(lengthBuffer >> 24),
			byte(lengthBuffer >> 16),
			byte(lengthBuffer >> 8),
			byte(lengthBuffer),
		})

		if err != nil {
			log.Println(err)
		}

		_, err = conn.Write(sendMessageBytes)
		if err != nil {
			log.Println(err)
		}
	}
}
