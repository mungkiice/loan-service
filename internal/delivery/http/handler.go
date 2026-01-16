package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mungkiice/-loan-service/internal/domain"
	"github.com/mungkiice/-loan-service/internal/usecase"
)

type Handler struct {
	loanUseCase *usecase.LoanUseCase
}

func NewHandler(loanUseCase *usecase.LoanUseCase) *Handler {
	return &Handler{loanUseCase: loanUseCase}
}

type CreateLoanRequest struct {
	BorrowerID      string  `json:"borrower_id" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required,gt=0"`
	Rate            float64 `json:"rate" binding:"required,gte=0"`
	ROI             float64 `json:"roi" binding:"required,gte=0"`
}

func (h *Handler) CreateLoan(c *gin.Context) {
	var req CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrowerID, err := uuid.Parse(req.BorrowerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid borrower_id"})
		return
	}

	loan, err := h.loanUseCase.CreateLoan(c.Request.Context(), usecase.CreateLoanRequest{
		BorrowerID:      borrowerID,
		PrincipalAmount: req.PrincipalAmount,
		Rate:            req.Rate,
		ROI:             req.ROI,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

type ApproveLoanRequest struct {
	ApprovalDate   string `form:"approval_date" binding:"required"`
	IdempotencyKey string `form:"idempotency_key" binding:"required"`
}

func (h *Handler) ApproveLoan(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}

	var req ApproveLoanRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eidStr, ok := userID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eid, err := uuid.Parse(eidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	approvalDate, err := parseTime(req.ApprovalDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad approval_date"})
		return
	}

	file, err := c.FormFile("picture_proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing picture_proof"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	if err := h.loanUseCase.ApproveLoan(c.Request.Context(), usecase.ApproveLoanRequest{
		LoanID:               loanID,
		EmployeeID:           eid,
		PictureProof:         f,
		PictureProofFilename: file.Filename,
		ApprovalDate:         approvalDate,
		IdempotencyKey:       req.IdempotencyKey,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type InvestRequest struct {
	InvestorID     string  `json:"investor_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	IdempotencyKey string  `json:"idempotency_key" binding:"required"`
}

func (h *Handler) Invest(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}

	var req InvestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	investorID, err := uuid.Parse(req.InvestorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investor_id"})
		return
	}

	if err := h.loanUseCase.Invest(c.Request.Context(), usecase.InvestRequest{
		LoanID:         loanID,
		InvestorID:     investorID,
		Amount:         req.Amount,
		IdempotencyKey: req.IdempotencyKey,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "investment added successfully"})
}

type DisburseLoanRequest struct {
	DisbursementDate string `form:"disbursement_date" binding:"required"`
	IdempotencyKey   string `form:"idempotency_key" binding:"required"`
}

func (h *Handler) DisburseLoan(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}

	var req DisburseLoanRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eidStr, ok := userID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	eid, err := uuid.Parse(eidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	disbursementDate, err := parseTime(req.DisbursementDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad disbursement_date"})
		return
	}

	file, err := c.FormFile("signed_agreement")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing signed_agreement"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	if err := h.loanUseCase.DisburseLoan(c.Request.Context(), usecase.DisburseLoanRequest{
		LoanID:                  loanID,
		EmployeeID:              eid,
		SignedAgreement:         f,
		SignedAgreementFilename: file.Filename,
		DisbursementDate:        disbursementDate,
		IdempotencyKey:          req.IdempotencyKey,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) GetLoan(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}

	loan, err := h.loanUseCase.GetLoan(c.Request.Context(), loanID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loan)
}

func (h *Handler) GetLoans(c *gin.Context) {
	stateStr := c.Query("state")
	if stateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state query parameter is required"})
		return
	}

	state := domain.LoanState(stateStr)
	if !isValidState(state) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	loans, err := h.loanUseCase.GetLoansByState(c.Request.Context(), state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loans)
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func isValidState(state domain.LoanState) bool {
	return state == domain.StateProposed ||
		state == domain.StateApproved ||
		state == domain.StateInvested ||
		state == domain.StateDisbursed
}
