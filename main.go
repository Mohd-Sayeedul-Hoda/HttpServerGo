package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
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

	header := extractHeader(req)

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

func extractHeader(req []byte) map[string]interface{} {
	header := make(map[string]interface{})
	header["headers"] = make(map[string]string)
	line := 0
	for i := 0; i < len(req); i++ {
		if req[i] == '\r' && req[i+1] == '\n' {
			line++
		}
	}
	line--
	headerLine := bytes.SplitN(req, []byte{'\r', '\n'}, line)
	for i, lines := range headerLine {
		if i == 0 {
			parts := bytes.Split(lines, []byte{' '})
			header["method"] = string(parts[0])
			header["url"] = string(parts[1])
			versionStr := string(parts[2])
			versionParts := strings.Split(versionStr, "/")
			if len(versionParts) == 2 && strings.HasPrefix(versionParts[1], "1.") {
				header["version"] = map[string]int{
					"major": 1,
					"minor": int(versionParts[1][2] - '0'),
				}
			}
		} else {
			parts := bytes.Split(lines, []byte{':'})
			fmt.Println(parts)
			key := strings.ToLower(strings.TrimSpace(string(parts[0])))
			value := strings.TrimSpace(string(parts[1]))
			header["headers"].(map[string]string)[key] = value
		}
	}
	fmt.Println(header)
	return header
}
