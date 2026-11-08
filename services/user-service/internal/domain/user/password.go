package user

// PasswordHash - Value Object for already calculated password hash
type PasswordHash struct {
	value string
}

// NewPasswordHash accepts only not empty hash
func NewPasswordHash(value string) (PasswordHash, error) {
	if value == "" {
		return PasswordHash{}, ErrInvalidPasswordHash
	}
	return PasswordHash{value: value}, nil
}

func (p PasswordHash) String() string {
	return p.value
}
