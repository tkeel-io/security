package license

type Licenser interface {
	Issue() ([]byte, error)
}

type License struct {
}

func NewLicense() Licenser {
	// TODO
	return License{}
}

func (l License) Issue() ([]byte, error) {
	// TODO
	return nil, nil
}
