package header

import (
	"bytes"
	"strings"
)

var Header map[string]interface{}

func ParseHeader(req []byte) map[string]interface{} {
	header := make(map[string]interface{})

	//  NOTE: declaring map inside map
	header["headers"] = make(map[string]string)

	//  NOTE: checking number of line of easily spliting the bytes
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

			//  NOTE: deriving the version of http header
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
			key := strings.ToLower(strings.TrimSpace(string(parts[0])))
			value := strings.TrimSpace(string(parts[1]))
			header["headers"].(map[string]string)[key] = value
		}
	}
	return header
}
