package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

type application struct {
	dir string
}

func main() {
	app := application{}
	flag.StringVar(&app.dir, "directory", "/tmp", "to specify the dirctory for file handler")
	err := app.Handler()
	if err != nil {
		log.Panic("Cannot register route", err)
	}
	flag.Parse()
	listen, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Panic(err)
		return
	}
	fmt.Println("Listning on port 8000")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("cannot accept the connection %s\n", err)
			return
		}
		go app.processRequest(conn)
	}
}
