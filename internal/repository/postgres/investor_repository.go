package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mungkiice/-loan-service/internal/domain"
)

// InvestorRepository implements domain.InvestorRepository using PostgreSQL
type InvestorRepository struct {
	db *pgxpool.Pool
}

// NewInvestorRepository creates a new investor repository
func NewInvestorRepository(db *pgxpool.Pool) *InvestorRepository {
	return &InvestorRepository{db: db}
}

// Create inserts a new investor
func (r *InvestorRepository) Create(ctx context.Context, investor *domain.Investor) error {
	query := `
		INSERT INTO investors (id, name, phone, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		investor.ID,
		investor.Name,
		investor.Phone,
		investor.Address,
		investor.CreatedAt,
		investor.UpdatedAt,
	)

	return err
}

// GetByID retrieves an investor by ID
func (r *InvestorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Investor, error) {
	query := `
		SELECT id, name, phone, address, created_at, updated_at
		FROM investors
		WHERE id = $1
	`

	var investor domain.Investor
	var phone, address sql.NullString

	err := r.db.QueryRow(ctx, query, id).Scan(
		&investor.ID,
		&investor.Name,
		&phone,
		&address,
		&investor.CreatedAt,
		&investor.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("investor not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	if phone.Valid {
		investor.Phone = &phone.String
	}
	if address.Valid {
		investor.Address = &address.String
	}

	investor.UserID = investor.ID // ID is the same as UserID in the new schema

	return &investor, nil
}

// GetByUserID retrieves an investor by user ID
func (r *InvestorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Investor, error) {
	query := `
		SELECT id, name, phone, address, created_at, updated_at
		FROM investors
		WHERE id = $1
	`

	var investor domain.Investor
	var phone, address sql.NullString

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&investor.ID,
		&investor.Name,
		&phone,
		&address,
		&investor.CreatedAt,
		&investor.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("investor not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	if phone.Valid {
		investor.Phone = &phone.String
	}
	if address.Valid {
		investor.Address = &address.String
	}

	investor.UserID = investor.ID

	return &investor, nil
}

// GetAll retrieves all investors
func (r *InvestorRepository) GetAll(ctx context.Context) ([]*domain.Investor, error) {
	query := `
		SELECT id, name, phone, address, created_at, updated_at
		FROM investors
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investors []*domain.Investor
	for rows.Next() {
		var investor domain.Investor
		var phone, address sql.NullString

		if err := rows.Scan(
			&investor.ID,
			&investor.Name,
			&phone,
			&address,
			&investor.CreatedAt,
			&investor.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if phone.Valid {
			investor.Phone = &phone.String
		}
		if address.Valid {
			investor.Address = &address.String
		}

		investor.UserID = investor.ID
		investors = append(investors, &investor)
	}

	return investors, rows.Err()
}
