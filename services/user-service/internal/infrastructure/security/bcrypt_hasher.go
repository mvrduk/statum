package security

import "context"
import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(_ context.Context, password string) (string, error) {
	value, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(value), nil
}
