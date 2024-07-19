package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Errorf("Cannot Listen to the port %s", err)
		return
	}
	fmt.Println("Listning on port 8000")
	_, err = listen.Accept()
	if err != nil {
		fmt.Errorf("cannot accept the connection %s", err)
		return
	}
	fmt.Println("connection accept")
}
