package main

import (
	"github.com/zipkero/ggnet/pkg/server"
	"log"
)

func main() {
	gameServer, err := server.NewServer("127.0.0.1:5000")
	if err != nil {
		log.Println(err)
	}
	err = gameServer.Start()
	if err != nil {
		log.Println(err)
	}
}
