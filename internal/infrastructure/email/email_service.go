package email

import (
	"context"
	"log"
	"os"
)

type EmailService interface {
	SendAgreementEmail(ctx context.Context, investorEmail string, agreementURL string) error
}

type MockEmailService struct {
	logger *log.Logger
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{
		logger: log.New(os.Stdout, "[EMAIL] ", log.LstdFlags),
	}
}

func (s *MockEmailService) SendAgreementEmail(ctx context.Context, investorEmail string, agreementURL string) error {
	s.logger.Printf("Sending agreement email to %s with URL: %s", investorEmail, agreementURL)
	return nil
}
