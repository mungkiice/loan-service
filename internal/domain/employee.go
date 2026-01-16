package domain

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Role      EmployeeRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EmployeeRole string

const (
	RoleFieldValidator EmployeeRole = "field_validator"
	RoleFieldOfficer   EmployeeRole = "field_officer"
	RoleAdmin          EmployeeRole = "admin"
)

func (r EmployeeRole) IsValid() bool {
	return r == RoleFieldValidator || r == RoleFieldOfficer || r == RoleAdmin
}

func NewEmployee(userID uuid.UUID, name string, role EmployeeRole) *Employee {
	now := time.Now()
	return &Employee{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
