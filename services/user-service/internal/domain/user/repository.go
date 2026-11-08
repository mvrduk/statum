package user

import "context"

// UserRepository contract for save and get user aggregate

type UserRepository interface {
	Create(ctx context.Context, u User) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
}
