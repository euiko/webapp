package webapp

type (
	Unmarshaler interface {
		Unmarshal(dst any) error
	}

	UnmarshalerFunc func(dst any) error
)

func (f UnmarshalerFunc) Unmarshal(dst any) error {
	return f(dst)
}
