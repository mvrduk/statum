package user

import (
	"context"
	"testing"
	"time"

	domain "github.com/mvrduk/statum/services/user-service/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) Create(ctx context.Context, u domain.User) (domain.User, error) {
	args := m.Called(ctx, u)

	if value := args.Get(0); value != nil {
		return value.(domain.User), args.Error(1)
	}
	return domain.User{}, args.Error(1)
}

func (m *repositoryMock) GetByID(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(ctx, id)
	if value := args.Get(0); value != nil {
		return value.(domain.User), args.Error(1)
	}
	return domain.User{}, args.Error(1)
}

type hasherStub struct{}

func (hasherStub) Hash(_ context.Context, password string) (string, error) {
	return "bcrypt-hash-of-" + password, nil
}

func TestRegister(t *testing.T) {
	repo := new(repositoryMock)
	svc := NewUserService(repo, hasherStub{})
	email, err := domain.NewEmail("test@example.com")
	assert.NoError(t, err)
	hash, err := domain.NewPasswordHash("bcrypt-hash-of-password123")
	assert.NoError(t, err)
	aggregate, err := domain.ReconstituteUser("user-123", email, hash, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Email().String() == "test@example.com" && u.PasswordHash().String() != "password123"
	})).Return(aggregate, nil).Once()

	got, err := svc.Register(context.Background(), RegisterInput{Email: "test@example.com", Password: "password123"})
	assert.NoError(t, err)
	assert.Equal(t, "user-123", got.ID)
	assert.Equal(t, "test@example.com", got.Email)
	repo.AssertExpectations(t)
}

func TestRegisterInvalidEmailDoesNotCallRepository(t *testing.T) {
	repo := new(repositoryMock)
	svc := NewUserService(repo, hasherStub{})
	_, err := svc.Register(context.Background(), RegisterInput{Email: "not-an-email", Password: "password123"})
	assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegisterShortPasswordDoesNotCallRepository(t *testing.T) {
	repo := new(repositoryMock)
	svc := NewUserService(repo, hasherStub{})
	_, err := svc.Register(context.Background(), RegisterInput{Email: "test@example.com", Password: "123"})
	assert.ErrorIs(t, err, domain.ErrInvalidPassword)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
