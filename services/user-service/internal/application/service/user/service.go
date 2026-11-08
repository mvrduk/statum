package user

import (
	"context"

	domain "github.com/mvrduk/statum/services/user-service/internal/domain/user"
)

type Service interface {
	Register(ctx context.Context, input RegisterInput) (UserDTO, error)
	GetByID(ctx context.Context, id string) (UserDTO, error)
}

type UserService struct {
	repo   domain.UserRepository
	hasher PasswordHasher
}

func NewUserService(repo domain.UserRepository, hasher PasswordHasher) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
	}
}

var _ Service = (*UserService)(nil)
