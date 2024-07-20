package main

import (
	"fmt"
	"github.com/zipkero/ggnet/pkg/server"
	"log"
)

func main() {
	gameServer, err := server.NewServer("127.0.0.1:5000")
	if err != nil {
		log.Println(err)
	}
	go func() {
		err = gameServer.Start()
		if err != nil {
			log.Println(err)
		}
	}()

	var cmd string
	for {
		_, err = fmt.Scanln(&cmd)
		if err != nil {
			log.Println(err)
		}
		switch cmd {
		case "exit":
			return
		case "kick":
			var sessionId string
			_, err = fmt.Scanln(&sessionId)
			if err != nil {
				log.Println(err)
				continue
			}
			err = gameServer.Kick(sessionId)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
