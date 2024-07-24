package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strconv"

	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/header"
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/request"
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/route"
)

func (app *application) processRequest(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)
	delimiter := []byte("\r\n\r\n")
	var reqHeader []byte
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			log.Printf("Error reading from connection: %s\n", err)
			return
		}
		reqHeader = append(reqHeader, line...)
		if bytes.HasSuffix(reqHeader, delimiter) {
			break
		}
	}

	Parseheader := header.ParseHeader(reqHeader)
	var body []byte
	if _, exists := Parseheader.Headers["content-length"]; exists {
		contLen, err := strconv.Atoi(Parseheader.Headers["content-length"])
		if err != nil {
			log.Println("cannot convert content lenght to int ", err)
			return
		}

		body = make([]byte, contLen)
		_, err = reader.Read(body)
		if err != nil {
			log.Println("Cannot read body from the reader ", err)
			return
		}
	}

	var resp string
	fn, param, err := route.ResolveRoute(Parseheader.URL, Parseheader.Method)
	if err != nil {
		header.NotFoundResp(conn)
		conn.Write([]byte(resp))
		return
	} else {
		Parseheader.Headers["accept-encoding"] = header.EcodingParser(Parseheader.Headers["accept-encoding"])
		req := request.Request{
			Header: Parseheader,
			Params: param,
			Body:   body,
		}
		err := fn(req, conn)
		if err != nil {
			log.Printf("Cannot write to connection %s", err)
			return
		}
	}
}
