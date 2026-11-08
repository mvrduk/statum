package user

import (
	"net/mail"
	"strings"
)

// Email - Value Object
type Email struct {
	value string
}

// NewEmail - create correct email adress
func NewEmail(raw string) (Email, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	addr, err := mail.ParseAddress(value)
	if err != nil || addr.Address != value {
		return Email{value: value}, ErrInvalidEmail
	}
	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) IsZero() bool {
	return e.value == ""
}
