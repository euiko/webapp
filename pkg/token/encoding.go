package token

type (
	Encoding interface {
		Encode(key Key, subject string, audiences ...string) ([]byte, error)
		Decode(key Key, b []byte) (*Token, error)
	}
)
