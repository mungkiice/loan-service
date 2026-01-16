package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LoanState string

const (
	StateProposed  LoanState = "proposed"
	StateApproved  LoanState = "approved"
	StateInvested  LoanState = "invested"
	StateDisbursed LoanState = "disbursed"
)

type Loan struct {
	ID                 uuid.UUID
	BorrowerID         uuid.UUID
	PrincipalAmount    float64
	Rate               float64
	ROI                float64
	AgreementLetterURL *string
	State              LoanState
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type LoanApproval struct {
	LoanID       uuid.UUID
	EmployeeID   uuid.UUID
	PictureProof string
	ApprovalDate time.Time
	CreatedAt    time.Time
}

type Investment struct {
	ID         uuid.UUID
	LoanID     uuid.UUID
	InvestorID uuid.UUID
	Amount     float64
	CreatedAt  time.Time
}

type Disbursement struct {
	LoanID             uuid.UUID
	EmployeeID         uuid.UUID
	SignedAgreementURL string
	DisbursementDate   time.Time
	CreatedAt          time.Time
}

type StateTransitionError struct {
	From  LoanState
	To    LoanState
	Cause string
}

var validTransitions = map[LoanState][]LoanState{
	StateProposed:  {StateApproved},
	StateApproved:  {StateInvested},
	StateInvested:  {StateDisbursed},
	StateDisbursed: {},
}

func (e *StateTransitionError) Error() string {
	return fmt.Sprintf("invalid transition %s -> %s: %s", e.From, e.To, e.Cause)
}

func (l *Loan) CanTransitionTo(target LoanState) error {
	allowed, exists := validTransitions[l.State]
	if !exists {
		return &StateTransitionError{
			From:  l.State,
			To:    target,
			Cause: "unknown current state",
		}
	}

	for _, allowedState := range allowed {
		if allowedState == target {
			return nil
		}
	}

	return &StateTransitionError{
		From:  l.State,
		To:    target,
		Cause: "transition not allowed",
	}
}

func (l *Loan) TransitionTo(newState LoanState) error {
	if err := l.CanTransitionTo(newState); err != nil {
		return err
	}

	l.State = newState
	l.UpdatedAt = time.Now()
	return nil
}

func NewLoan(borrowerID uuid.UUID, principalAmount, rate, roi float64) *Loan {
	now := time.Now()
	return &Loan{
		ID:              uuid.New(),
		BorrowerID:      borrowerID,
		PrincipalAmount: principalAmount,
		Rate:            rate,
		ROI:             roi,
		State:           StateProposed,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (l *Loan) ValidateInvestmentAmount(amount float64, currentTotal float64) error {
	if amount <= 0 {
		return errors.New("investment amount must be positive")
	}

	if currentTotal+amount > l.PrincipalAmount {
		return fmt.Errorf("total investment (%.2f) would exceed principal (%.2f)", currentTotal+amount, l.PrincipalAmount)
	}

	return nil
}

func (l *Loan) IsFullyInvested(totalInvested float64) bool {
	const epsilon = 0.01
	return totalInvested >= l.PrincipalAmount-epsilon
}
