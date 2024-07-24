package main

import (
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/header"
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/request"
)

func (app *application) EchoHandler(req request.Request, conn net.Conn) error {
	respEncode, err := header.Encode([]byte(req.Params["params"]), req.Header.Headers["accept-encoding"])
	if err != nil {
		err = header.InternalServerErrorResp(conn)
		return err
	}
	contentLength := strconv.Itoa(len(respEncode))
	writeHead := header.RespHeader{
		VersionMajor:    req.Header.Version["major"],
		VersionMinor:    req.Header.Version["minor"],
		ContentType:     "text/plain",
		Status:          200,
		ContentLenght:   contentLength,
		ContentEncoding: req.Header.Headers["accept-encoding"],
	}
	err = writeHead.WriteHeader(conn)
	if err != nil {
		return err
	}

	_, err = conn.Write(respEncode)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) RootHandler(req request.Request, conn net.Conn) error {
	writeHead := header.RespHeader{
		VersionMajor:    req.Header.Version["major"],
		VersionMinor:    req.Header.Version["minor"],
		Status:          200,
		ContentEncoding: req.Header.Headers["accept-encoding"],
	}
	err := writeHead.WriteHeader(conn)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) UserAgentHandler(req request.Request, conn net.Conn) error {

	respEncode, err := header.Encode([]byte(req.Header.Headers["user-agent"]), req.Header.Headers["accept-encoding"])
	if err != nil {
		err = header.InternalServerErrorResp(conn)
		return err
	}
	contentLength := strconv.Itoa(len(respEncode))
	writeHead := header.RespHeader{
		VersionMajor:    req.Header.Version["major"],
		VersionMinor:    req.Header.Version["minor"],
		Status:          200,
		ContentType:     "text/plain",
		ContentLenght:   contentLength,
		ContentEncoding: req.Header.Headers["accept-encoding"],
	}
	err = writeHead.WriteHeader(conn)
	if err != nil {
		return err
	}
	_, err = conn.Write(respEncode)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) FileHandler(req request.Request, conn net.Conn) error {
	reqFile := req.Params["params"]
	filePath := filepath.Join(app.dir, reqFile)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = header.NotFoundResp(conn)
		return err
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		err = header.InternalServerErrorResp(conn)
		return err
	}

	respEncode, err := header.Encode(data, req.Header.Headers["accept-encoding"])
	if err != nil {
		err = header.InternalServerErrorResp(conn)
		return err
	}

	contentLen := strconv.Itoa(len(respEncode))

	writeHeader := header.RespHeader{
		VersionMajor:    req.Header.Version["major"],
		VersionMinor:    req.Header.Version["minor"],
		Status:          200,
		ContentType:     "application/octet-stream",
		ContentLenght:   contentLen,
		ContentEncoding: req.Header.Headers["accept-encoding"],
	}

	err = writeHeader.WriteHeader(conn)
	if err != nil {
		return err
	}

	// Write file content
	_, err = conn.Write(respEncode)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) FileCreateHandler(req request.Request, conn net.Conn) error {
	fileName := req.Params["params"]
	filePath := filepath.Join(app.dir, fileName)

	// Create or overwrite the file
	err := os.WriteFile(filePath, req.Body, 0644)
	if err != nil {
		err = header.InternalServerErrorResp(conn)
		return err
	}

	// Send success response
	writeHeader := header.RespHeader{
		VersionMajor:    req.Header.Version["major"],
		VersionMinor:    req.Header.Version["minor"],
		Status:          201,
		ContentEncoding: req.Header.Headers["content-encoding"],
	}
	err = writeHeader.WriteHeader(conn)
	if err != nil {
		return err
	}

	return nil
}
