package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mungkiice/-loan-service/internal/domain"
)

// ApprovalRepository implements domain.ApprovalRepository using PostgreSQL
type ApprovalRepository struct {
	db *pgxpool.Pool
}

// NewApprovalRepository creates a new approval repository
func NewApprovalRepository(db *pgxpool.Pool) *ApprovalRepository {
	return &ApprovalRepository{db: db}
}

// Create inserts a new approval
func (r *ApprovalRepository) Create(ctx context.Context, approval *domain.LoanApproval) error {
	query := `
		INSERT INTO loan_approvals (loan_id, employee_id, picture_proof, approval_date, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		approval.LoanID,
		approval.EmployeeID,
		approval.PictureProof,
		approval.ApprovalDate,
		approval.CreatedAt,
	)

	return err
}

// GetByLoanID retrieves approval by loan ID
func (r *ApprovalRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.LoanApproval, error) {
	query := `
		SELECT loan_id, employee_id, picture_proof, approval_date, created_at
		FROM loan_approvals
		WHERE loan_id = $1
	`

	var approval domain.LoanApproval
	err := r.db.QueryRow(ctx, query, loanID).Scan(
		&approval.LoanID,
		&approval.EmployeeID,
		&approval.PictureProof,
		&approval.ApprovalDate,
		&approval.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("approval not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	return &approval, nil
}
