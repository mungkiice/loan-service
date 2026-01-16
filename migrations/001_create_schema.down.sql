-- Drop triggers
DROP TRIGGER IF EXISTS update_loans_updated_at ON loans;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS disbursements;
DROP TABLE IF EXISTS investments;
DROP TABLE IF EXISTS loan_approvals;
DROP TABLE IF EXISTS loans;

-- Drop enum type
DROP TYPE IF EXISTS loan_state;

-- Drop triggers
DROP TRIGGER IF EXISTS update_investors_updated_at ON investors;
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS investors CASCADE;
DROP TABLE IF EXISTS employees CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop enum types
DROP TYPE IF EXISTS employee_role;
DROP TYPE IF EXISTS user_type;
