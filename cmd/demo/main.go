package main

import (
	"context"

	"github.com/euiko/webapp"
)

func main() {
	app := webapp.New("go-fullstack-boilerplate", "WEBAPP")

	// Service modules
	app.Register(newHelloService)
	app.Run(context.Background())
}
