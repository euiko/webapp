package session

type (
	Encoding interface {
		Decode() (*Session, error)
		Encode(*Session) error
	}
)
