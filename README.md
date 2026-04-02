# Finance Data Processing and Access Control Backend

A Go backend for a finance dashboard system with role-based access control, financial records management, and analytics APIs.

## Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin (HTTP router)
- **Database:** PostgreSQL
- **Authentication:** JWT (HS256)
- **Password Hashing:** bcrypt (cost 12)

## Project Structure

```
finance-backend/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── config/                     # Environment-based configuration
│   ├── domain/                     # Entities, interfaces, domain errors
│   ├── handler/                    # HTTP handlers (thin layer)
│   ├── middleware/                 # JWT auth + RBAC middleware
│   ├── repository/postgres/        # Database implementations
│   ├── service/                    # Business logic layer
│   └── router/                    # Route definitions + middleware wiring
├── pkg/pagination/                 # Shared pagination utilities
├── migrations/                     # SQL migration files
├── Makefile
└── README.md
```

**Architecture:** 3-layer clean architecture (Handler → Service → Repository). Domain layer defines interfaces; repository layer implements them. No business logic in handlers, no HTTP concerns in services.

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 13+

### 1. Create the database

```bash
createdb finance_db
```

### 2. Set environment variables

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_pg_user
export DB_PASSWORD=your_pg_password
export DB_NAME=finance_db
export DB_SSLMODE=disable
export JWT_SECRET=your-secret-key-change-this
export SERVER_PORT=8080
```

### 3. Run the server

```bash
make run
```

Migrations run automatically on startup.

### 4. Run tests

```bash
make test
```

## API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register a new user (default role: viewer) |
| POST | `/api/v1/auth/login` | Login and receive JWT token |

### User Profile (Any authenticated user)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get own profile |
| PUT | `/api/v1/users/me` | Update own profile (name, password) |

### User Management (Admin only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List all users (paginated) |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user (name, role, status) |
| DELETE | `/api/v1/users/:id` | Soft delete user |

### Financial Records (Admin: full CRUD, Analyst: read-only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/transactions` | Create transaction (Admin) |
| GET | `/api/v1/transactions` | List transactions with filters (Analyst, Admin) |
| GET | `/api/v1/transactions/:id` | Get transaction by ID (Analyst, Admin) |
| PUT | `/api/v1/transactions/:id` | Update transaction (Admin) |
| DELETE | `/api/v1/transactions/:id` | Soft delete transaction (Admin) |

**Query parameters for listing:**
- `type` — filter by `income` or `expense`
- `category` — filter by category name
- `date_from` — filter from date (YYYY-MM-DD)
- `date_to` — filter to date (YYYY-MM-DD)
- `sort_by` — `date`, `amount`, `created_at`, `category` (default: `date`)
- `sort_order` — `asc` or `desc` (default: `desc`)
- `page` — page number (default: 1)
- `page_size` — items per page (default: 20, max: 100)

### Dashboard (Viewer, Analyst, Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/dashboard/summary` | Total income, expenses, net balance |
| GET | `/api/v1/dashboard/category-totals` | Totals grouped by category and type |
| GET | `/api/v1/dashboard/trends` | Monthly or weekly income/expense trends |
| GET | `/api/v1/dashboard/recent` | Recent transactions |

**Query parameters:**
- `date_from`, `date_to` — date range filter
- `period` — `monthly` or `weekly` (trends endpoint)
- `limit` — number of recent transactions (default: 10)

## Roles and Access Control

| Role | Dashboard | View Records | Create/Edit/Delete Records | Manage Users |
|------|-----------|-------------|---------------------------|-------------|
| **Viewer** | Yes | No | No | No |
| **Analyst** | Yes | Yes | No | No |
| **Admin** | Yes | Yes | Yes | Yes |

Access control is enforced via middleware. The JWT auth middleware validates tokens and checks that the user is still active. The RBAC middleware checks the user's role against the required roles for each endpoint.

## Authentication

All protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

Register a user and receive a token:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "password123", "name": "Admin User"}'
```

Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "password123"}'
```

## Error Response Format

All errors follow a consistent structure:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {"field": "amount", "message": "must be greater than 0"}
    ]
  }
}
```

**HTTP status codes used:**
- `400` — Validation error
- `401` — Unauthorized (missing/invalid token)
- `403` — Forbidden (insufficient role)
- `404` — Resource not found
- `409` — Conflict (e.g., duplicate email)
- `500` — Internal server error

## Design Decisions and Assumptions

1. **Soft delete** — Financial records and users are never permanently deleted. A `deleted_at` timestamp is set instead, and all queries filter out soft-deleted rows.

2. **Transactions are global** — The `user_id` field tracks who created the record, not ownership. All authorized users see all transactions (company ledger model).

3. **Roles are fixed** — Three roles (viewer, analyst, admin) stored as a CHECK-constrained string column. No separate roles table since the set is small and static.

4. **New users default to viewer role** — An admin must explicitly promote users to analyst or admin.

5. **Single JWT with 24h TTL** — No refresh token flow. The auth middleware checks `is_active` on every request, so deactivating a user effectively revokes access.

6. **Raw SQL over ORM** — Dashboard aggregate queries (conditional SUM, date_trunc) are more natural in SQL. Two tables don't justify the ORM abstraction cost.

7. **Amount stored as positive number** — The `type` field (income/expense) carries the sign semantics. `NUMERIC(15,2)` in PostgreSQL avoids floating-point precision issues.

8. **Migrations run on startup** — Simple migration runner tracks applied versions. Suitable for development; production deployments may prefer a separate migration step.
