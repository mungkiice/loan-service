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

type LoanRepository struct {
	db *pgxpool.Pool
}

func NewLoanRepository(db *pgxpool.Pool) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	query := `
		INSERT INTO loans (id, borrower_id, principal_amount, rate, roi, agreement_letter_url, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(ctx, query,
		loan.ID,
		loan.BorrowerID,
		loan.PrincipalAmount,
		loan.Rate,
		loan.ROI,
		loan.AgreementLetterURL,
		loan.State,
		loan.CreatedAt,
		loan.UpdatedAt,
	)

	return err
}

func (r *LoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	query := `
		SELECT id, borrower_id, principal_amount, rate, roi, agreement_letter_url, state, created_at, updated_at
		FROM loans
		WHERE id = $1
	`

	var loan domain.Loan
	var agreementLetterURL sql.NullString

	err := r.db.QueryRow(ctx, query, id).Scan(
		&loan.ID,
		&loan.BorrowerID,
		&loan.PrincipalAmount,
		&loan.Rate,
		&loan.ROI,
		&agreementLetterURL,
		&loan.State,
		&loan.CreatedAt,
		&loan.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("loan not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	if agreementLetterURL.Valid {
		loan.AgreementLetterURL = &agreementLetterURL.String
	}

	return &loan, nil
}

func (r *LoanRepository) GetByState(ctx context.Context, state domain.LoanState) ([]*domain.Loan, error) {
	query := `
		SELECT id, borrower_id, principal_amount, rate, roi, agreement_letter_url, state, created_at, updated_at
		FROM loans
		WHERE state = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := make([]*domain.Loan, 0)
	for rows.Next() {
		var loan domain.Loan
		var agreementLetterURL sql.NullString

		if err := rows.Scan(
			&loan.ID,
			&loan.BorrowerID,
			&loan.PrincipalAmount,
			&loan.Rate,
			&loan.ROI,
			&agreementLetterURL,
			&loan.State,
			&loan.CreatedAt,
			&loan.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if agreementLetterURL.Valid {
			loan.AgreementLetterURL = &agreementLetterURL.String
		}

		loans = append(loans, &loan)
	}

	return loans, rows.Err()
}

func (r *LoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	query := `
		UPDATE loans
		SET principal_amount = $2, rate = $3, roi = $4, agreement_letter_url = $5, state = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		loan.ID,
		loan.PrincipalAmount,
		loan.Rate,
		loan.ROI,
		loan.AgreementLetterURL,
		loan.State,
		loan.UpdatedAt,
	)

	return err
}
