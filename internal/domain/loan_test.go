package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewLoan(t *testing.T) {
	borrowerID := uuid.New()
	loan := NewLoan(borrowerID, 10000.0, 5.0, 3.0)

	assert.NotNil(t, loan)
	assert.Equal(t, borrowerID, loan.BorrowerID)
	assert.Equal(t, 10000.0, loan.PrincipalAmount)
	assert.Equal(t, 5.0, loan.Rate)
	assert.Equal(t, 3.0, loan.ROI)
	assert.Equal(t, StateProposed, loan.State)
	assert.NotEqual(t, uuid.Nil, loan.ID)
}

func TestCanTransitionTo(t *testing.T) {
	tests := []struct {
		name      string
		from      LoanState
		to        LoanState
		shouldErr bool
	}{
		{"proposed to approved", StateProposed, StateApproved, false},
		{"approved to invested", StateApproved, StateInvested, false},
		{"invested to disbursed", StateInvested, StateDisbursed, false},
		{"proposed to invested", StateProposed, StateInvested, true},
		{"approved to proposed", StateApproved, StateProposed, true},
		{"disbursed to any", StateDisbursed, StateApproved, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loan := &Loan{State: tt.from}
			err := loan.CanTransitionTo(tt.to)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransitionTo(t *testing.T) {
	loan := &Loan{State: StateProposed}

	err := loan.TransitionTo(StateApproved)
	assert.NoError(t, err)
	assert.Equal(t, StateApproved, loan.State)

	err = loan.TransitionTo(StateInvested)
	assert.NoError(t, err)
	assert.Equal(t, StateInvested, loan.State)

	err = loan.TransitionTo(StateDisbursed)
	assert.NoError(t, err)
	assert.Equal(t, StateDisbursed, loan.State)
}

func TestValidateInvestmentAmount(t *testing.T) {
	loan := &Loan{PrincipalAmount: 10000.0}

	tests := []struct {
		name        string
		amount      float64
		currentTotal float64
		shouldErr   bool
	}{
		{"valid amount", 5000.0, 0.0, false},
		{"exact principal", 10000.0, 0.0, false},
		{"exceeds principal", 5001.0, 5000.0, true},
		{"zero amount", 0.0, 0.0, true},
		{"negative amount", -100.0, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loan.ValidateInvestmentAmount(tt.amount, tt.currentTotal)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsFullyInvested(t *testing.T) {
	loan := &Loan{PrincipalAmount: 10000.0}

	assert.True(t, loan.IsFullyInvested(10000.0))
	assert.True(t, loan.IsFullyInvested(10000.01))
	assert.False(t, loan.IsFullyInvested(9999.0))
	assert.False(t, loan.IsFullyInvested(5000.0))
}
