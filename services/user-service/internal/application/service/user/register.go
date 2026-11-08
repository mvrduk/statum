package user

import (
	"context"

	domain "github.com/mvrduk/statum/services/user-service/internal/domain/user"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 72
)

func (s *UserService) Register(ctx context.Context, input RegisterInput) (UserDTO, error) {
	if len(input.Password) < minPasswordLength || len(input.Password) > maxPasswordLength {
		return UserDTO{}, domain.ErrInvalidPassword
	}

	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return UserDTO{}, err
	}

	hash, err := s.hasher.Hash(ctx, input.Password)
	if err != nil {
		return UserDTO{}, err
	}

	passwordHash, err := domain.NewPasswordHash(hash)
	if err != nil {
		return UserDTO{}, err
	}

	aggregate, err := domain.NewUser(email, passwordHash)
	if err != nil {
		return UserDTO{}, err
	}

	saved, err := s.repo.Create(ctx, aggregate)
	if err != nil {
		return UserDTO{}, err
	}

	return toDTO(saved), nil
}
