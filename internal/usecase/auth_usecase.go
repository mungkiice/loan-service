package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mungkiice/-loan-service/internal/domain"
	"github.com/mungkiice/-loan-service/internal/infrastructure/jwt"
)

type AuthUseCase struct {
	userRepo     domain.UserRepository
	employeeRepo domain.EmployeeRepository
	investorRepo domain.InvestorRepository
	jwtService   *jwt.JWTService
}

// NewAuthUseCase creates a new auth use case
// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	userRepo domain.UserRepository,
	employeeRepo domain.EmployeeRepository,
	investorRepo domain.InvestorRepository,
	jwtService *jwt.JWTService,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
		investorRepo: investorRepo,
		jwtService:   jwtService,
	}
}

type SignInRequest struct {
	Email    string
	Password string
}

type SignInResponse struct {
	Token     string
	User      map[string]interface{}
	ExpiresIn int64
}

// SignIn authenticates a user and returns a JWT token
func (uc *AuthUseCase) SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !domain.CheckPassword(user.Password, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	var role string
	var data map[string]interface{}

	switch user.UserType {
	case domain.UserTypeEmployee:
		emp, err := uc.employeeRepo.GetByUserID(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("employee not found")
		}
		role = string(emp.Role)
		data = map[string]interface{}{
			"id":    emp.ID.String(),
			"name":  emp.Name,
			"role":  role,
			"email": user.Email,
		}
	case domain.UserTypeInvestor:
		inv, err := uc.investorRepo.GetByUserID(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("investor not found")
		}
		data = map[string]interface{}{
			"id":    inv.ID.String(),
			"name":  inv.Name,
			"email": user.Email,
		}
	default:
		return nil, fmt.Errorf("unknown user type")
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, string(user.UserType), role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &SignInResponse{
		Token:     token,
		User:      data,
		ExpiresIn: int64(uc.jwtService.TokenDuration().Seconds()),
	}, nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	claims, err := uc.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	if _, err := uc.userRepo.GetByID(ctx, claims.UserID); err != nil {
		return nil, errors.New("user not found")
	}

	return claims, nil
}
