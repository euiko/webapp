package auth

type (
	UserLoader[User any] interface {
		LoadUser(loginId string, password string) (*User, error)
	}
)
