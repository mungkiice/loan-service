#  Loan Service

A production-grade Loan Engine service built with Go, implementing Clean Architecture principles. The service manages loan lifecycle with strict state transitions and supports multiple investors.

## Loan States

The loan lifecycle follows these states:

1. **proposed** (initial state)
   - Loan is created and awaiting approval

2. **approved**
   - Requires: picture proof, employee_id, approval date
   - Cannot return to proposed state

3. **invested**
   - Multiple investors can invest
   - Total investment cannot exceed principal
   - When fully invested, automatically transitions to invested state
   - All investors receive agreement letter email

4. **disbursed** (terminal state)
   - Requires: signed agreement letter, employee_id, disbursement date
   - Final state, no further transitions allowed

## State Transition Rules

- State transitions can only move forward (no backward transitions)
- Each transition has specific requirements that must be met
- All transitions are atomic and transactional
- Idempotency keys prevent duplicate operations

## API Endpoints

### Loan Lifecycle

#### Create Loan
```http
POST /api/v1/loans
Content-Type: application/json

{
  "borrower_id": "uuid",
  "principal_amount": 10000.00,
  "rate": 5.0,
  "roi": 3.0
}
```

#### Approve Loan
```http
POST /api/v1/loans/{id}/approve
Content-Type: multipart/form-data

employee_id: uuid
approval_date: 2024-01-01T00:00:00Z (RFC3339)
idempotency_key: unique-key
picture_proof: <file>
```

#### Invest in Loan
```http
POST /api/v1/loans/{id}/invest
Content-Type: application/json

{
  "investor_id": "uuid",
  "amount": 5000.00,
  "idempotency_key": "unique-key"
}
```

#### Disburse Loan
```http
POST /api/v1/loans/{id}/disburse
Content-Type: multipart/form-data

employee_id: uuid
disbursement_date: 2024-01-01T00:00:00Z (RFC3339)
idempotency_key: unique-key
signed_agreement: <file>
```

#### Get Loan
```http
GET /api/v1/loans/{id}
```

#### Get Loans by State
```http
GET /api/v1/loans?state=proposed
GET /api/v1/loans?state=approved
GET /api/v1/loans?state=invested
GET /api/v1/loans?state=disbursed
```

## Database Schema

### Tables

- **loans**: Main loan entity
- **loan_approvals**: Approval information
- **investments**: Investment records (multiple per loan)
- **disbursements**: Disbursement information

All tables include proper indexing, foreign keys, and constraints.

## Running Locally

### Prerequisites

- Go
- PostgreSQL
- Redis

### Manual Setup

1. Set up PostgreSQL and Redis

2. Set environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/loan_db?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export FILE_STORAGE_PATH="./storage"
export FILE_STORAGE_URL="http://localhost:8080/files"
export PORT="8080"
```

3. Run migrations:
```bash
psql -U postgres -d loan_db -f migrations/001_create_schema.up.sql
```

4. Build and run:
```bash
go mod download
go build -o loan-service ./cmd/api
./loan-service
```

## Testing

### Unit Tests

Run unit tests:
```bash
go test ./internal/domain/...
go test ./internal/usecase/...
```

### All Tests

Run all tests with coverage:
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Example API Requests

### 1. Create a Loan

```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": "550e8400-e29b-41d4-a716-446655440000",
    "principal_amount": 10000.00,
    "rate": 5.0,
    "roi": 3.0
  }'
```

### 2. Approve the Loan

```bash
curl -X POST http://localhost:8080/api/v1/loans/{loan_id}/approve \
  -F "employee_id=550e8400-e29b-41d4-a716-446655440001" \
  -F "approval_date=2024-01-01T00:00:00Z" \
  -F "idempotency_key=approve-001" \
  -F "picture_proof=@proof.jpg"
```

### 3. Invest in the Loan

```bash
curl -X POST http://localhost:8080/api/v1/loans/{loan_id}/invest \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "550e8400-e29b-41d4-a716-446655440002",
    "amount": 5000.00,
    "idempotency_key": "invest-001"
  }'
```

### 4. Complete Investment (when total reaches principal)

```bash
# Make another investment to reach full amount
curl -X POST http://localhost:8080/api/v1/loans/{loan_id}/invest \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "550e8400-e29b-41d4-a716-446655440003",
    "amount": 5000.00,
    "idempotency_key": "invest-002"
  }'
```

### 5. Disburse the Loan

```bash
curl -X POST http://localhost:8080/api/v1/loans/{loan_id}/disburse \
  -F "employee_id=550e8400-e29b-41d4-a716-446655440004" \
  -F "disbursement_date=2024-01-02T00:00:00Z" \
  -F "idempotency_key=disburse-001" \
  -F "signed_agreement=@agreement.pdf"
```

### 6. Get Loan Details

```bash
curl http://localhost:8080/api/v1/loans/{loan_id}
```

### 7. Get Loans by State

```bash
curl http://localhost:8080/api/v1/loans?state=approved
```

## Business Rules

1. **State Transitions**:
   - Cannot approve unless in `proposed` state
   - Cannot invest unless in `approved` state
   - Cannot disburse unless in `invested` state
   - All transitions are forward-only

2. **Investments**:
   - Total investments must not exceed principal amount
   - When total investment equals principal, loan automatically transitions to `invested`
   - All investors receive email with agreement letter URL when fully invested

3. **Idempotency**:
   - All state transition operations require idempotency keys
   - Duplicate requests with same key are rejected

4. **Concurrency**:
   - Investment operations use Redis locks to prevent race conditions
   - Database transactions ensure data consistency