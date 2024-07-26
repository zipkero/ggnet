package main

import (
	"fmt"
	"github.com/zipkero/ggnet/pkg/client"
	"github.com/zipkero/ggnet/pkg/message"
	"log"
)

func main() {
	// create client
	c := client.NewClient("127.0.0.1", "5000")
	if err := c.Connect(); err != nil {
		log.Fatalln(err)
	}

	// receive message
	go func() {
		for {
			select {
			case msg := <-c.ReceiveCh:
				fmt.Println(msg.Content)
			}
		}
	}()

	// send message
	go func() {
		sendMessage := message.Message{
			Type: 1,
		}
		for {
			_, err := fmt.Scanln(&sendMessage.Content)
			if err != nil {
				log.Println(err)
			}
			c.SendCh <- sendMessage
		}
	}()

	// listen
	if err := c.Listen(); err != nil {
		log.Fatalln(err)
	}
}
