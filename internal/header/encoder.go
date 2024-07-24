package header

import (
	"bytes"
	"compress/gzip"
	"strings"
)

var supportedEncoding = "gzip"

func EcodingParser(encode string) string {
	encodings := strings.Split(encode, ",")
	for _, value := range encodings {
		if strings.TrimSpace(value) == supportedEncoding {
			return value
		}
	}
	return ""
}

func Encode(data []byte, algo string) ([]byte, error) {
	if algo == supportedEncoding {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err := zw.Write(data)
		if err != nil {
			return nil, err
		}
		zw.Close()

		// Compressed data
		data = buf.Bytes()
	}
	return data, nil
}
