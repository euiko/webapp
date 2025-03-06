package main

import (
	"context"
	"embed"
	"errors"
	"log"

	"github.com/euiko/webapp"
	"github.com/euiko/webapp/db/sqldb"
	"github.com/euiko/webapp/module/auth"
	"github.com/euiko/webapp/module/static"
)

var (
	//go:embed db/migrations
	migrations embed.FS
)

type (
	User struct {
		LoginId  string `json:"loginId"`
		Password string `json:"password"`
	}
	userLoader struct{}
)

var (
	demoUser = User{
		LoginId:  "demo",
		Password: "12345678",
	}
)

func (u User) Subject() string {
	return u.LoginId
}

func newUserLoader() *userLoader {
	return &userLoader{}
}

func (l *userLoader) LoadUser(loginId string, password string) (*User, error) {
	if loginId == demoUser.LoginId && password == demoUser.Password {
		return &User{
			LoginId:  demoUser.LoginId,
			Password: demoUser.Password,
		}, nil
	}

	return nil, errors.New("invalid login")
}

func main() {
	// register migrations
	sqldb.AddMigrationFS(migrations)

	app := webapp.New("demo", "WEBAPP_DEMO")

	// Service modules
	app.Register(static.ModuleFactory())
	app.Register(auth.ModuleFactory(newUserLoader()))
	app.Register(newHelloService)
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
