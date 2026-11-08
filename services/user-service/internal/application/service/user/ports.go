package user

import "context"

type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
}
