package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/header"
	"log"
	"net"
)

func main() {
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
		go resp(conn)
	}
}

func resp(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)
	delimiter := []byte("\r\n\r\n")
	var req []byte
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			log.Printf("Error reading from connection: %s\n", err)
			return
		}
		req = append(req, line...)
		if bytes.HasSuffix(req, delimiter) {
			break
		}
	}

	header := header.ParseHeader(req)

	var resp string
	if header["url"] == "/" {
		resp = fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		resp = fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	_, err := conn.Write([]byte(resp))
	if err != nil {
		log.Printf("Cannot write to the connection %s\n", err)
		return
	}
}
