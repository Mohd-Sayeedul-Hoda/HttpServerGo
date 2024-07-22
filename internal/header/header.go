package header

import (
	"bytes"
	"strings"
)

type Header struct {
	Method  string
	URL     string
	Version map[string]int
	Headers map[string]string
}

func ParseHeader(req []byte) Header {
	header := Header{
		Headers: make(map[string]string),
		Version: make(map[string]int),
	}

	line := 0
	for i := 0; i < len(req); i++ {
		if req[i] == '\r' && i+1 < len(req) && req[i+1] == '\n' {
			line++
		}
	}
	line--

	headerLines := bytes.SplitN(req, []byte{'\r', '\n'}, line)
	for i, lines := range headerLines {
		if i == 0 {
			parts := bytes.Split(lines, []byte{' '})
			if len(parts) >= 3 {
				header.Method = string(parts[0])
				header.URL = string(parts[1])
				versionStr := string(parts[2])
				versionParts := strings.Split(versionStr, "/")
				if len(versionParts) == 2 && strings.HasPrefix(versionParts[1], "1.") {
					header.Version["major"] = 1
					header.Version["minor"] = int(versionParts[1][2] - '0')
				}
			}
		} else {
			parts := bytes.SplitN(lines, []byte{':'}, 2)
			if len(parts) == 2 {
				key := strings.ToLower(strings.TrimSpace(string(parts[0])))
				value := strings.TrimSpace(string(parts[1]))
				header.Headers[key] = value
			}
		}
	}

	return header
}
