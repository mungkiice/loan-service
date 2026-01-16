package domain

import (
	"context"

	"github.com/google/uuid"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *Loan) error
	GetByID(ctx context.Context, id uuid.UUID) (*Loan, error)
	GetByState(ctx context.Context, state LoanState) ([]*Loan, error)
	Update(ctx context.Context, loan *Loan) error
}

type ApprovalRepository interface {
	Create(ctx context.Context, approval *LoanApproval) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) (*LoanApproval, error)
}

type InvestmentRepository interface {
	Create(ctx context.Context, investment *Investment) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) ([]*Investment, error)
	GetTotalByLoanID(ctx context.Context, loanID uuid.UUID) (float64, error)
}

type DisbursementRepository interface {
	Create(ctx context.Context, disbursement *Disbursement) error
	GetByLoanID(ctx context.Context, loanID uuid.UUID) (*Disbursement, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type EmployeeRepository interface {
	Create(ctx context.Context, employee *Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*Employee, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Employee, error)
	GetAll(ctx context.Context) ([]*Employee, error)
}

type InvestorRepository interface {
	Create(ctx context.Context, investor *Investor) error
	GetByID(ctx context.Context, id uuid.UUID) (*Investor, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Investor, error)
	GetAll(ctx context.Context) ([]*Investor, error)
}
