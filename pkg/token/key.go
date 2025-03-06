package token

type (
	Key interface {
		Public() []byte
		Private() []byte
	}

	SymetricKey []byte

	AsymetricKey struct {
		public  []byte
		private []byte
	}
)

func NewSymetricKey(key []byte) SymetricKey {
	return key
}

func NewAsymetricKey(public, private []byte) AsymetricKey {
	return AsymetricKey{
		public:  public,
		private: private,
	}
}

func (k SymetricKey) Public() []byte {
	return k
}

func (k SymetricKey) Private() []byte {
	return k
}

func (k AsymetricKey) Public() []byte {
	return k.public
}

func (k AsymetricKey) Private() []byte {
	return k.private
}
