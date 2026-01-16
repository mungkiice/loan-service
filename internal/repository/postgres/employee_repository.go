package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mungkiice/-loan-service/internal/domain"
)

// EmployeeRepository implements domain.EmployeeRepository using PostgreSQL
type EmployeeRepository struct {
	db *pgxpool.Pool
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create inserts a new employee
func (r *EmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	query := `
		INSERT INTO employees (id, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		employee.ID,
		employee.Name,
		employee.Role,
		employee.CreatedAt,
		employee.UpdatedAt,
	)

	return err
}

// GetByID retrieves an employee by ID
func (r *EmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT id, name, role, created_at, updated_at
		FROM employees
		WHERE id = $1
	`

	var employee domain.Employee
	err := r.db.QueryRow(ctx, query, id).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Role,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	employee.UserID = employee.ID // ID is the same as UserID in the new schema

	return &employee, nil
}

// GetByUserID retrieves an employee by user ID
func (r *EmployeeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT id, name, role, created_at, updated_at
		FROM employees
		WHERE id = $1
	`

	var employee domain.Employee
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Role,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	employee.UserID = employee.ID

	return &employee, nil
}

// GetAll retrieves all employees
func (r *EmployeeRepository) GetAll(ctx context.Context) ([]*domain.Employee, error) {
	query := `
		SELECT id, name, role, created_at, updated_at
		FROM employees
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*domain.Employee
	for rows.Next() {
		var employee domain.Employee
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Role,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		); err != nil {
			return nil, err
		}
		employee.UserID = employee.ID
		employees = append(employees, &employee)
	}

	return employees, rows.Err()
}
