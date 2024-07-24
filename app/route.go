package main

import (
	"github.com/Mohd-Sayeedul-Hoda/httpServer/internal/route"
)

func (app *application) Handler() error {
	err := route.RegisterRoute("GET", "/", app.RootHandler)
	if err != nil {
		return err
	}
	err = route.RegisterRoute("GET", "/echo/:param", app.EchoHandler)
	if err != nil {
		return err
	}
	err = route.RegisterRoute("GET", "/user-agent", app.UserAgentHandler)
	if err != nil {
		return err
	}

	err = route.RegisterRoute("GET", "/files/:param", app.FileHandler)
	if err != nil {
		return err
	}

	err = route.RegisterRoute("POST", "/files/:param", app.FileCreateHandler)
	if err != nil {
		return err
	}
	return nil
}
