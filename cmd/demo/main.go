package main

import (
	"context"
	"embed"
	"log"

	"github.com/euiko/webapp"
	"github.com/euiko/webapp/db/sqldb"
	"github.com/euiko/webapp/module/static"
)

var (
	//go:embed db/migrations
	migrations embed.FS
)

func main() {
	// register migrations
	sqldb.AddMigrationFS(migrations)

	app := webapp.New("demo", "WEBAPP_DEMO")

	// Service modules
	app.Register(static.NewModule)
	app.Register(newHelloService)
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
