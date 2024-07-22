package route

import (
	"errors"
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/request"
	"net"
	"strings"
	"sync"
)

type fn func(request.Request, net.Conn) error

var (
	routes = make(map[string]fn)
	mu     sync.Mutex
)

func RegisterRoute(path string, fn func(request.Request, net.Conn) error) error {
	mu.Lock()
	defer mu.Unlock()
	_, exist := routes[path]
	if exist {
		return errors.New("URL already exists")
	}
	routes[path] = fn
	return nil
}

func ResolveRoute(path string) (fn, map[string]string, error) {
	mu.Lock()
	defer mu.Unlock()
	for route, handler := range routes {
		if params, match := matchRoute(route, path); match {
			return handler, params, nil
		}
	}
	return nil, nil, errors.New("URL does't exist")
}

func matchRoute(route, path string) (map[string]string, bool) {
	routeSegments := strings.Split(route, "/")
	pathSegments := strings.Split(path, "/")
	if len(routeSegments) != len(pathSegments) {
		return nil, false
	}

	params := make(map[string]string)
	for i, routeSegment := range routeSegments {
		if strings.HasPrefix(routeSegment, ":") {
			params["params"] = pathSegments[i]
		} else if routeSegment == pathSegments[i] {
			continue
		} else {
			return nil, false
		}
	}

	return params, true
}
