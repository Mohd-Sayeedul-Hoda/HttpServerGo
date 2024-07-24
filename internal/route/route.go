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
	routesGET  = make(map[string]fn)
	routesPOST = make(map[string]fn)
	mu         sync.Mutex
)

func RegisterRoute(method string, path string, fn func(request.Request, net.Conn) error) error {
	mu.Lock()
	defer mu.Unlock()
	if method == "GET" {
		_, exist := routesGET[path]
		if exist {
			return errors.New("URL already exists")
		}
		routesGET[path] = fn
	} else if method == "POST" {
		_, exist := routesPOST[path]
		if exist {
			return errors.New("URL already exists")
		}
		routesPOST[path] = fn
	}
	return nil
}

func ResolveRoute(path string, method string) (fn, map[string]string, error) {
	mu.Lock()
	defer mu.Unlock()
	if method == "GET" {
		for route, handler := range routesGET {
			if params, match := matchRoute(route, path); match {
				return handler, params, nil
			}
		}
	} else if method == "POST" {
		for route, handler := range routesPOST {
			if params, match := matchRoute(route, path); match {
				return handler, params, nil
			}
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
