package header

import (
	"fmt"
	"net"
)

type RespHeader struct {
	ContentLenght   string
	ContentType     string
	ContentEncoding string
	Status          Status
	VersionMajor    int
	VersionMinor    int
}

type Status int

const (
	OK                  Status = 200
	Created             Status = 201
	NotFound            Status = 404
	InternalServerError Status = 500
)

func (headResp *RespHeader) WriteHeader(conn net.Conn) error {

	resp := fmt.Sprintf("HTTP/%d.%d %d %s\r\n", headResp.VersionMajor, headResp.VersionMinor, headResp.Status, statusToMessage(headResp.Status))
	if headResp.ContentEncoding != "" {
		con := fmt.Sprintf("Content-Encoding: %s\r\n", headResp.ContentEncoding)
		resp += con
	}
	if headResp.ContentType != "" {
		con := fmt.Sprintf("Content-Type: %s\r\n", headResp.ContentType)
		resp += con
	}
	if headResp.ContentLenght != "" {
		con := fmt.Sprintf("Content-Length: %s\r\n", headResp.ContentLenght)
		resp += con
	}
	resp += "\r\n"
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func NotFoundResp(conn net.Conn) error {
	resp := "HTTP/1.1 404 Not Found\r\n\r\n"
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func InternalServerErrorResp(conn net.Conn) error {
	resp := "HTTP/1.1 500 Internal Server Error\r\n\r\n"
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}
	return nil
}

func statusToMessage(status Status) string {
	var ret string
	switch status {
	case OK:
		ret = "OK"
	case Created:
		ret = "Created"
	case NotFound:
		ret = "Not Found"
	case InternalServerError:
		ret = "Internal Server Error"
	default:
		ret = ""
	}
	return ret
}
