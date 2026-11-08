package user

import "time"

type RegisterInput struct {
	Email    string
	Password string
}

type UserDTO struct {
	ID        string
	Email     string
	CreatedAt time.Time
}
