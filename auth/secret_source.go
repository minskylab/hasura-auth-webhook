package auth

// SecretSource ...
type SecretSource interface {
	GetAccessSecret() ([]byte, error)
	GetRefreshSecret() ([]byte, error)
}

type rawSecret struct {
	accessSecret  []byte
	refreshSecret []byte
}

// GetSecret ...
func (s rawSecret) GetAccessSecret() ([]byte, error) {
	return s.accessSecret, nil
}

// GetSecret ...
func (s rawSecret) GetRefreshSecret() ([]byte, error) {
	return s.refreshSecret, nil
}
