package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mungkiice/-loan-service/internal/domain"
	"github.com/mungkiice/-loan-service/internal/infrastructure/email"
	"github.com/mungkiice/-loan-service/internal/infrastructure/redis"
	"github.com/mungkiice/-loan-service/internal/infrastructure/storage"
)

type LoanUseCase struct {
	loanRepo         domain.LoanRepository
	approvalRepo     domain.ApprovalRepository
	investmentRepo   domain.InvestmentRepository
	disbursementRepo domain.DisbursementRepository
	userRepo         domain.UserRepository
	redisClient      redis.RedisClient
	fileStorage      storage.FileStorage
	emailService     email.EmailService
}

func NewLoanUseCase(
	loanRepo domain.LoanRepository,
	approvalRepo domain.ApprovalRepository,
	investmentRepo domain.InvestmentRepository,
	disbursementRepo domain.DisbursementRepository,
	userRepo domain.UserRepository,
	redisClient redis.RedisClient,
	fileStorage storage.FileStorage,
	emailService email.EmailService,
) *LoanUseCase {
	return &LoanUseCase{
		loanRepo:         loanRepo,
		approvalRepo:     approvalRepo,
		investmentRepo:   investmentRepo,
		disbursementRepo: disbursementRepo,
		userRepo:         userRepo,
		redisClient:      redisClient,
		fileStorage:      fileStorage,
		emailService:     emailService,
	}
}

func (uc *LoanUseCase) CreateLoan(ctx context.Context, req CreateLoanRequest) (*domain.Loan, error) {
	loan := domain.NewLoan(
		req.BorrowerID,
		req.PrincipalAmount,
		req.Rate,
		req.ROI,
	)

	if err := uc.loanRepo.Create(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	return loan, nil
}

func (uc *LoanUseCase) ApproveLoan(ctx context.Context, req ApproveLoanRequest) error {
	idempotencyKey := fmt.Sprintf("approve:%s:%s", req.LoanID, req.IdempotencyKey)
	if exists, _ := uc.redisClient.CheckIdempotencyKey(ctx, idempotencyKey); exists {
		return fmt.Errorf("duplicate request: idempotency key already used")
	}

	loan, err := uc.loanRepo.GetByID(ctx, req.LoanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}

	if err := loan.CanTransitionTo(domain.StateApproved); err != nil {
		return err
	}

	picturePath, err := uc.fileStorage.Store(ctx, req.PictureProof, req.PictureProofFilename)
	if err != nil {
		return fmt.Errorf("failed to store picture proof: %w", err)
	}

	approval := &domain.LoanApproval{
		LoanID:       req.LoanID,
		EmployeeID:   req.EmployeeID,
		PictureProof: uc.fileStorage.GetURL(picturePath),
		ApprovalDate: req.ApprovalDate,
		CreatedAt:    time.Now(),
	}

	if err := loan.TransitionTo(domain.StateApproved); err != nil {
		return err
	}

	if err := uc.loanRepo.Update(ctx, loan); err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	if err := uc.approvalRepo.Create(ctx, approval); err != nil {
		return fmt.Errorf("failed to create approval: %w", err)
	}

	_ = uc.redisClient.SetIdempotencyKey(ctx, idempotencyKey, "approved", 24*time.Hour)

	return nil
}

func (uc *LoanUseCase) Invest(ctx context.Context, req InvestRequest) error {
	lockKey := fmt.Sprintf("invest:%s", req.LoanID)
	acquired, err := uc.redisClient.AcquireLock(ctx, lockKey, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !acquired {
		return fmt.Errorf("could not acquire lock, please try again")
	}
	defer uc.redisClient.ReleaseLock(ctx, lockKey)

	idempotencyKey := fmt.Sprintf("invest:%s:%s:%s", req.LoanID, req.InvestorID, req.IdempotencyKey)
	if exists, _ := uc.redisClient.CheckIdempotencyKey(ctx, idempotencyKey); exists {
		return fmt.Errorf("duplicate request: idempotency key already used")
	}

	loan, err := uc.loanRepo.GetByID(ctx, req.LoanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}

	if loan.State != domain.StateApproved {
		return fmt.Errorf("loan must be in approved state to accept investments")
	}

	currentTotal, err := uc.investmentRepo.GetTotalByLoanID(ctx, req.LoanID)
	if err != nil {
		return fmt.Errorf("failed to get current investment total: %w", err)
	}

	if err := loan.ValidateInvestmentAmount(req.Amount, currentTotal); err != nil {
		return err
	}

	investment := &domain.Investment{
		ID:         uuid.New(),
		LoanID:     req.LoanID,
		InvestorID: req.InvestorID,
		Amount:     req.Amount,
		CreatedAt:  time.Now(),
	}

	if err := uc.investmentRepo.Create(ctx, investment); err != nil {
		return fmt.Errorf("failed to create investment: %w", err)
	}

	newTotal := currentTotal + req.Amount
	if loan.IsFullyInvested(newTotal) {
		if err := loan.TransitionTo(domain.StateInvested); err != nil {
			return err
		}

		if err := uc.loanRepo.Update(ctx, loan); err != nil {
			return fmt.Errorf("failed to update loan state: %w", err)
		}

		agreementURL := uc.fileStorage.GetURL(fmt.Sprintf("agreements/%s.pdf", req.LoanID))
		loan.AgreementLetterURL = &agreementURL
		if err := uc.loanRepo.Update(ctx, loan); err != nil {
			return fmt.Errorf("failed to update agreement letter URL: %w", err)
		}

		investments, err := uc.investmentRepo.GetByLoanID(ctx, req.LoanID)
		if err == nil {
			for _, inv := range investments {
				investorUser, err := uc.userRepo.GetByID(ctx, inv.InvestorID)
				if err == nil {
					_ = uc.emailService.SendAgreementEmail(ctx, investorUser.Email, agreementURL)
				}
			}
		}
	}

	_ = uc.redisClient.SetIdempotencyKey(ctx, idempotencyKey, "invested", 24*time.Hour)

	return nil
}

func (uc *LoanUseCase) DisburseLoan(ctx context.Context, req DisburseLoanRequest) error {
	idempotencyKey := fmt.Sprintf("disburse:%s:%s", req.LoanID, req.IdempotencyKey)
	if exists, _ := uc.redisClient.CheckIdempotencyKey(ctx, idempotencyKey); exists {
		return fmt.Errorf("duplicate request: idempotency key already used")
	}

	loan, err := uc.loanRepo.GetByID(ctx, req.LoanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}

	if err := loan.CanTransitionTo(domain.StateDisbursed); err != nil {
		return err
	}

	agreementPath, err := uc.fileStorage.Store(ctx, req.SignedAgreement, req.SignedAgreementFilename)
	if err != nil {
		return fmt.Errorf("failed to store signed agreement: %w", err)
	}

	disbursement := &domain.Disbursement{
		LoanID:             req.LoanID,
		EmployeeID:         req.EmployeeID,
		SignedAgreementURL: uc.fileStorage.GetURL(agreementPath),
		DisbursementDate:   req.DisbursementDate,
		CreatedAt:          time.Now(),
	}

	if err := loan.TransitionTo(domain.StateDisbursed); err != nil {
		return err
	}

	if err := uc.loanRepo.Update(ctx, loan); err != nil {
		return fmt.Errorf("failed to update loan: %w", err)
	}

	if err := uc.disbursementRepo.Create(ctx, disbursement); err != nil {
		return fmt.Errorf("failed to create disbursement: %w", err)
	}

	_ = uc.redisClient.SetIdempotencyKey(ctx, idempotencyKey, "disbursed", 24*time.Hour)

	return nil
}

func (uc *LoanUseCase) GetLoan(ctx context.Context, loanID uuid.UUID) (*domain.Loan, error) {
	loan, err := uc.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("loan:%s", loanID)
	_ = uc.redisClient.SetCache(ctx, cacheKey, loanID.String(), 5*time.Minute)

	return loan, nil
}

func (uc *LoanUseCase) GetLoansByState(ctx context.Context, state domain.LoanState) ([]*domain.Loan, error) {
	return uc.loanRepo.GetByState(ctx, state)
}

type CreateLoanRequest struct {
	BorrowerID      uuid.UUID
	PrincipalAmount float64
	Rate            float64
	ROI             float64
}

type ApproveLoanRequest struct {
	LoanID               uuid.UUID
	EmployeeID           uuid.UUID
	PictureProof         interface{ Read([]byte) (int, error) }
	PictureProofFilename string
	ApprovalDate         time.Time
	IdempotencyKey       string
}

type InvestRequest struct {
	LoanID         uuid.UUID
	InvestorID     uuid.UUID
	Amount         float64
	IdempotencyKey string
}

type DisburseLoanRequest struct {
	LoanID                  uuid.UUID
	EmployeeID              uuid.UUID
	SignedAgreement         interface{ Read([]byte) (int, error) }
	SignedAgreementFilename string
	DisbursementDate        time.Time
	IdempotencyKey          string
}
