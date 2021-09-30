package auth

// SecretSource ...
type SecretSource interface {
	GetSecret() ([]byte, error)
}

type rawSecret struct {
	secret []byte
}

// GetSecret ...
func (s rawSecret) GetSecret() ([]byte, error) {
	return s.secret, nil
}
