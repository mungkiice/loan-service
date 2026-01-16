package usecase

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mungkiice/-loan-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock repositories
type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) Create(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetByState(ctx context.Context, state domain.LoanState) ([]*domain.Loan, error) {
	args := m.Called(ctx, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Loan), args.Error(1)
}

func (m *MockLoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

type MockApprovalRepository struct {
	mock.Mock
}

func (m *MockApprovalRepository) Create(ctx context.Context, approval *domain.LoanApproval) error {
	args := m.Called(ctx, approval)
	return args.Error(0)
}

func (m *MockApprovalRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.LoanApproval, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoanApproval), args.Error(1)
}

type MockInvestmentRepository struct {
	mock.Mock
}

func (m *MockInvestmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockInvestmentRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) ([]*domain.Investment, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetTotalByLoanID(ctx context.Context, loanID uuid.UUID) (float64, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(float64), args.Error(1)
}

type MockDisbursementRepository struct {
	mock.Mock
}

func (m *MockDisbursementRepository) Create(ctx context.Context, disbursement *domain.Disbursement) error {
	args := m.Called(ctx, disbursement)
	return args.Error(0)
}

func (m *MockDisbursementRepository) GetByLoanID(ctx context.Context, loanID uuid.UUID) (*domain.Disbursement, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Disbursement), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockRedisClient implements redis.RedisClient interface for testing
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) SetIdempotencyKey(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) CheckIdempotencyKey(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	args := m.Called(ctx, key, expiration)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) ReleaseLock(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRedisClient) SetCache(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) GetCache(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) Store(ctx context.Context, file io.Reader, filename string) (string, error) {
	args := m.Called(ctx, file, filename)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorage) GetURL(path string) string {
	args := m.Called(path)
	return args.String(0)
}

func (m *MockFileStorage) Delete(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendAgreementEmail(ctx context.Context, investorEmail string, agreementURL string) error {
	args := m.Called(ctx, investorEmail, agreementURL)
	return args.Error(0)
}

func TestCreateLoan(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockApprovalRepository)
	mockInvestmentRepo := new(MockInvestmentRepository)
	mockDisbursementRepo := new(MockDisbursementRepository)
	mockUserRepo := new(MockUserRepository)
	mockRedis := new(MockRedisClient)
	mockFileStorage := new(MockFileStorage)
	mockEmail := new(MockEmailService)

	uc := NewLoanUseCase(
		mockLoanRepo,
		mockApprovalRepo,
		mockInvestmentRepo,
		mockDisbursementRepo,
		mockUserRepo,
		mockRedis,
		mockFileStorage,
		mockEmail,
	)

	borrowerID := uuid.New()
	req := CreateLoanRequest{
		BorrowerID:      borrowerID,
		PrincipalAmount: 10000.0,
		Rate:            5.0,
		ROI:             3.0,
	}

	mockLoanRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)

	loan, err := uc.CreateLoan(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, loan)
	assert.Equal(t, domain.StateProposed, loan.State)
	mockLoanRepo.AssertExpectations(t)
}

func TestApproveLoan(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockApprovalRepository)
	mockInvestmentRepo := new(MockInvestmentRepository)
	mockDisbursementRepo := new(MockDisbursementRepository)
	mockUserRepo := new(MockUserRepository)
	mockRedis := new(MockRedisClient)
	mockFileStorage := new(MockFileStorage)
	mockEmail := new(MockEmailService)

	uc := NewLoanUseCase(
		mockLoanRepo,
		mockApprovalRepo,
		mockInvestmentRepo,
		mockDisbursementRepo,
		mockUserRepo,
		mockRedis,
		mockFileStorage,
		mockEmail,
	)

	loanID := uuid.New()
	employeeID := uuid.New()
	loan := domain.NewLoan(uuid.New(), 10000.0, 5.0, 3.0)
	loan.ID = loanID

	mockRedis.On("CheckIdempotencyKey", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	mockLoanRepo.On("GetByID", mock.Anything, loanID).Return(loan, nil)
	mockFileStorage.On("Store", mock.Anything, mock.Anything, mock.AnythingOfType("string")).Return("proof.jpg", nil)
	mockFileStorage.On("GetURL", "proof.jpg").Return("http://example.com/proof.jpg")
	mockLoanRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Loan")).Return(nil)
	mockApprovalRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.LoanApproval")).Return(nil)
	mockRedis.On("SetIdempotencyKey", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(nil)

	req := ApproveLoanRequest{
		LoanID:               loanID,
		EmployeeID:           employeeID,
		PictureProof:         bytes.NewReader([]byte("fake image")),
		PictureProofFilename: "proof.jpg",
		ApprovalDate:         time.Now(),
		IdempotencyKey:       "test-key",
	}

	err := uc.ApproveLoan(context.Background(), req)

	require.NoError(t, err)
	mockLoanRepo.AssertExpectations(t)
	mockApprovalRepo.AssertExpectations(t)
}
