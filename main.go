package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/header"
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/request"
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/route"
)

func main() {
	err := route.RegisterRoute("/", rootHandler)
	if err != nil {
		log.Panic(err)
	}
	err = route.RegisterRoute("/echo/:param", echoHandler)
	if err != nil {
		log.Panic(err)
	}
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
	fn, param, err := route.ResolveRoute(header.URL)
	if err != nil {
		resp = fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
		conn.Write([]byte(resp))
		return
	} else {
		req := request.Request{
			Header: header,
			Params: param,
		}
		fn(req, conn)
	}
	_, err = conn.Write([]byte(resp))
	if err != nil {
		log.Printf("Cannot write to the connection %s\n", err)
		return
	}
}

func echoHandler(req request.Request, conn net.Conn) error {
	respString := req.Params["params"]
	contentLength := len(respString)
	resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, respString)
	conn.Write([]byte(resp))
	return nil
}

func rootHandler(req request.Request, conn net.Conn) error {
	resp := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n")
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}
