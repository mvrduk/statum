package user

import "time"

// User is Aggreate root for agregate of user
type User struct {
	id           string
	email        Email
	passwordHash PasswordHash
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUser creates new user for checked Value Objects
func NewUser(email Email, passwordHash PasswordHash) (User, error) {
	if email.IsZero() || passwordHash.value == "" {
		return User{}, ErrInvalidUser
	}
	return User{email: email, passwordHash: passwordHash}, nil
}

func ReconstituteUser(id string, email Email, passwordHash PasswordHash, createdAt, updatedAt time.Time) (User, error) {
	if id == "" || email.IsZero() || passwordHash.value == "" || createdAt.IsZero() {
		return User{}, ErrInvalidUser
	}
	return User{id: id, email: email, createdAt: createdAt, updatedAt: updatedAt}, nil
}

func (u User) ID() string                 { return u.id }
func (u User) Email() Email               { return u.email }
func (u User) PasswordHash() PasswordHash { return u.passwordHash }
func (u User) CreatedAt() time.Time       { return u.createdAt }
func (u User) UpdatedAt() time.Time       { return u.updatedAt }
