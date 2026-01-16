package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	UserTypeEmployee UserType = "employee"
	UserTypeInvestor UserType = "investor"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	UserType  UserType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Investor struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Phone     *string
	Address   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func NewInvestor(userID uuid.UUID, name string, phone, address *string) *Investor {
	now := time.Now()
	return &Investor{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Phone:     phone,
		Address:   address,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
