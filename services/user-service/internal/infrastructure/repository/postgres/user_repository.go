package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	domain "github.com/mvrduk/statum/services/user-service/internal/domain/user"
)

const (
	pgUniqueViolation = "23505"
	pgInvalidTextRepr = "22P02"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }

var _ domain.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	const query = `
			INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id::text, email, password_hash, created_at, updated_at
    `

	var id, email, passwordHash string
	var createdAt, updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, u.Email().String(), u.PasswordHash().String()).Scan(&id, &email, &passwordHash, &createdAt, &updatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.User{}, domain.ErrEmailAlreadyUsed
		}

		return domain.User{}, err
	}

	return reconstitute(id, email, passwordHash, createdAt.Time, updatedAt.Time)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	const query = `
		SELECT id::text, email, password_hash, created_at, updated_at
		FROM users WHERE id = $1
    `
	var UserID, email, passwordHash string
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(&UserID, &email, &passwordHash, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgInvalidTextRepr {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return reconstitute(UserID, email, passwordHash, createdAt.Time, updatedAt.Time)
}

func reconstitute(id, email, hash string, createdAt, updatedAt time.Time) (domain.User, error) {
	emailVO, err := domain.NewEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	if err != nil {
		return domain.User{}, err
	}
	hashVO, err := domain.NewPasswordHash(hash)
	if err != nil {
		return domain.User{}, err
	}
	return domain.ReconstituteUser(id, emailVO, hashVO, createdAt, updatedAt)
}
