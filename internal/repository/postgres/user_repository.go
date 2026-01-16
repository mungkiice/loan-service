package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mungkiice/-loan-service/internal/domain"
)

// UserRepository implements domain.UserRepository using PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password, user_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.UserType,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password, user_type, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, user_type, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
