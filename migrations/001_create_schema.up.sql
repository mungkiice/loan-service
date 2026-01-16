-- Create enum type for loan state
CREATE TYPE loan_state AS ENUM ('proposed', 'approved', 'invested', 'disbursed');

-- Create loans table
CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id UUID NOT NULL,
    principal_amount DECIMAL(15, 2) NOT NULL CHECK (principal_amount > 0),
    rate DECIMAL(5, 2) NOT NULL CHECK (rate >= 0),
    roi DECIMAL(5, 2) NOT NULL CHECK (roi >= 0),
    agreement_letter_url TEXT,
    state loan_state NOT NULL DEFAULT 'proposed',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on state for filtering
CREATE INDEX idx_loans_state ON loans(state);
CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_created_at ON loans(created_at);

-- Create loan_approvals table
CREATE TABLE loan_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL UNIQUE REFERENCES loans(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    picture_proof TEXT NOT NULL,
    approval_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_loan_approvals_loan_id ON loan_approvals(loan_id);
CREATE INDEX idx_loan_approvals_employee_id ON loan_approvals(employee_id);

-- Create investments table
CREATE TABLE investments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    investor_id UUID NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_investments_loan_id ON investments(loan_id);
CREATE INDEX idx_investments_investor_id ON investments(investor_id);

-- Create disbursements table
CREATE TABLE disbursements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL UNIQUE REFERENCES loans(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL,
    signed_agreement_url TEXT NOT NULL,
    disbursement_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_disbursements_loan_id ON disbursements(loan_id);
CREATE INDEX idx_disbursements_employee_id ON disbursements(employee_id);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_loans_updated_at BEFORE UPDATE ON loans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create user_type enum
CREATE TYPE user_type AS ENUM ('employee', 'investor');

-- Create users table (base table for authentication)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    user_type user_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on email for fast lookups
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_user_type ON users(user_type);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create employee_role enum
CREATE TYPE employee_role AS ENUM ('field_validator', 'field_officer', 'admin');

-- Create employees table with user reference
CREATE TABLE employees (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role employee_role NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_employees_role ON employees(role);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_employees_updated_at BEFORE UPDATE ON employees
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create investors table
CREATE TABLE investors (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_investors_name ON investors(name);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_investors_updated_at BEFORE UPDATE ON investors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Seed sample users and employees
-- Password for all: "password123" (hashed with bcrypt)
INSERT INTO users (id, email, password, user_type) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'validator1@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'employee'),
('550e8400-e29b-41d4-a716-446655440002', 'validator2@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'employee'),
('550e8400-e29b-41d4-a716-446655440003', 'officer1@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'employee'),
('550e8400-e29b-41d4-a716-446655440004', 'officer2@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'employee'),
('550e8400-e29b-41d4-a716-446655440005', 'admin@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'employee'),
('550e8400-e29b-41d4-a716-446655440010', 'investor1@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'investor'),
('550e8400-e29b-41d4-a716-446655440011', 'investor2@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'investor'),
('550e8400-e29b-41d4-a716-446655440012', 'investor3@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'investor'),
('550e8400-e29b-41d4-a716-446655440013', 'investor4@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'investor'),
('550e8400-e29b-41d4-a716-446655440014', 'investor5@.com', '$2a$10$gopcD59oMAj.03TuQbLshO7eHp.yqEqfzn1TbEsFLrAWfCxE.mTEO', 'investor');

-- Seed employees
INSERT INTO employees (id, name, role) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Field Validator 1', 'field_validator'),
('550e8400-e29b-41d4-a716-446655440002', 'Field Validator 2', 'field_validator'),
('550e8400-e29b-41d4-a716-446655440003', 'Field Officer 1', 'field_officer'),
('550e8400-e29b-41d4-a716-446655440004', 'Field Officer 2', 'field_officer'),
('550e8400-e29b-41d4-a716-446655440005', 'Admin User', 'admin');

-- Seed investors
INSERT INTO investors (id, name, phone, address) VALUES
('550e8400-e29b-41d4-a716-446655440010', 'Investor 1', '+6281234567890', 'Jakarta, Indonesia'),
('550e8400-e29b-41d4-a716-446655440011', 'Investor 2', '+6281234567891', 'Bandung, Indonesia'),
('550e8400-e29b-41d4-a716-446655440012', 'Investor 3', '+6281234567892', 'Surabaya, Indonesia'),
('550e8400-e29b-41d4-a716-446655440013', 'Investor 4', '+6281234567893', 'Yogyakarta, Indonesia'),
('550e8400-e29b-41d4-a716-446655440014', 'Investor 5', '+6281234567894', 'Medan, Indonesia');
