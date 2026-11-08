package user

import (
	"context"

	domain "github.com/mvrduk/statum/services/user-service/internal/domain/user"
)

func (s *UserService) GetByID(ctx context.Context, id string) (UserDTO, error) {
	aggregate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return UserDTO{}, err
	}

	return toDTO(aggregate), nil
}

func toDTO(u domain.User) UserDTO {
	return UserDTO{ID: u.ID(), Email: u.Email().String(), CreatedAt: u.CreatedAt()}
}
