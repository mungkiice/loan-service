package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mungkiice/-loan-service/internal/domain"
)

// DisbursementRepository implements domain.DisbursementRepository using PostgreSQL
type DisbursementRepository struct {
	db *pgxpool.Pool
}

// NewDisbursementRepository creates a new disbursement repository
func NewDisbursementRepository(db *pgxpool.Pool) *DisbursementRepository {
	return &DisbursementRepository{db: db}
}

// Create inserts a new disbursement
func (r *DisbursementRepository) Create(ctx context.Context, disbursement *domain.Disbursement) error {
	query := `
		INSERT INTO disbursements (loan_id, employee_id, signed_agreement_url, disbursement_date, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		disbursement.LoanID,
		disbursement.EmployeeID,
		disbursement.SignedAgreementURL,
		disbursement.DisbursementDate,
		disbursement.CreatedAt,
	)

	return err
}

// GetByLoanID retrieves disbursement by loan ID
func (r *DisbursementRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Disbursement, error) {
	query := `
		SELECT loan_id, employee_id, signed_agreement_url, disbursement_date, created_at
		FROM disbursements
		WHERE loan_id = $1
	`

	var disbursement domain.Disbursement
	err := r.db.QueryRow(ctx, query, loanID).Scan(
		&disbursement.LoanID,
		&disbursement.EmployeeID,
		&disbursement.SignedAgreementURL,
		&disbursement.DisbursementDate,
		&disbursement.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("disbursement not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	return &disbursement, nil
}
