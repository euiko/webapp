package main

import (
	"context"

	"github.com/euiko/webapp"
	"github.com/euiko/webapp/module/static"
)

func main() {
	app := webapp.New("go-fullstack-boilerplate", "WEBAPP")

	// Service modules
	app.Register(static.NewModule)
	app.Register(newHelloService)
	app.Run(context.Background())
}
