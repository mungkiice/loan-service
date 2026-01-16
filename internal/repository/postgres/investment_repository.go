package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mungkiice/-loan-service/internal/domain"
)

// InvestmentRepository implements domain.InvestmentRepository using PostgreSQL
type InvestmentRepository struct {
	db *pgxpool.Pool
}

// NewInvestmentRepository creates a new investment repository
func NewInvestmentRepository(db *pgxpool.Pool) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

// Create inserts a new investment
func (r *InvestmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	query := `
		INSERT INTO investments (id, loan_id, investor_id, amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		investment.ID,
		investment.LoanID,
		investment.InvestorID,
		investment.Amount,
		investment.CreatedAt,
	)

	return err
}

// GetByLoanID retrieves all investments for a loan
func (r *InvestmentRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) ([]*domain.Investment, error) {
	query := `
		SELECT id, loan_id, investor_id, amount, created_at
		FROM investments
		WHERE loan_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []*domain.Investment
	for rows.Next() {
		var investment domain.Investment
		if err := rows.Scan(
			&investment.ID,
			&investment.LoanID,
			&investment.InvestorID,
			&investment.Amount,
			&investment.CreatedAt,
		); err != nil {
			return nil, err
		}
		investments = append(investments, &investment)
	}

	return investments, rows.Err()
}

// GetTotalByLoanID calculates the total investment amount for a loan
func (r *InvestmentRepository) GetTotalByLoanID(ctx context.Context, loanID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM investments
		WHERE loan_id = $1
	`

	var total float64
	err := r.db.QueryRow(ctx, query, loanID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total investment: %w", err)
	}

	return total, nil
}
