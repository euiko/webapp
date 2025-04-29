package main

import (
	"context"
	"embed"
	"errors"
	"log"

	"github.com/euiko/webapp"
	"github.com/euiko/webapp/db/sqldb"
	"github.com/euiko/webapp/module/auth"
	"github.com/euiko/webapp/module/rbac"
	"github.com/euiko/webapp/module/static"
	"github.com/mitchellh/mapstructure"
)

var (
	//go:embed db/migrations
	migrations embed.FS
)

type (
	User struct {
		LoginId  string `json:"login_id"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	userSession struct {
		LoginId string `json:"login_id" mapstructure:"login_id"`
		Role    string `json:"role" mapstructure:"role"`
	}

	userLoader struct{}
)

var (
	demoUser = User{
		LoginId:  "demo",
		Password: "12345678",
		Role:     "admin",
	}
)

func (u *User) LoginID() string {
	return u.LoginId
}

func (u *User) Name() string {
	return u.LoginId
}

func (u *User) RoleName() string {
	return u.Role
}

func (u *User) TokenSubject() string {
	return u.LoginId
}

func (u *User) MarshalSession() (interface{}, error) {
	return userSession{
		LoginId: u.LoginId,
		Role:    u.Role,
	}, nil
}

func (u *User) UnmarshalSession(value interface{}) error {
	var userSession userSession
	if err := mapstructure.Decode(value, &userSession); err != nil {
		return err
	}

	u.LoginId = userSession.LoginId
	u.Role = userSession.Role
	return nil
}

func newUserLoader() *userLoader {
	return &userLoader{}
}

func (l *userLoader) UserById(ctx context.Context, loginId string) (*User, error) {
	if loginId == demoUser.LoginId {
		return &User{
			LoginId:  demoUser.LoginId,
			Password: demoUser.Password,
			Role:     demoUser.Role,
		}, nil
	}

	return nil, errors.New("users not found")
}

func (l *userLoader) LoadUser(ctx context.Context, loginId string, password string) (*User, error) {
	log.Println("loginId", loginId, "password", password)
	if loginId == demoUser.LoginId && password == demoUser.Password {
		return &User{
			LoginId:  demoUser.LoginId,
			Password: demoUser.Password,
			Role:     demoUser.Role,
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
	app.Register(rbac.ModuleFactory())
	app.Register(newHelloService)
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
