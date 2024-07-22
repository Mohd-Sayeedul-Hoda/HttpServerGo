package request

import (
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/header"
)

type Request struct {
	Header header.Header
	Params map[string]string
}
