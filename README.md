# 🚗 SpotSync Backend API

Smart parking and EV spot reservation backend built for the B6A6 assignment.

![Go](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go&logoColor=white)
![Echo](https://img.shields.io/badge/Echo-v5-222222)
![GORM](https://img.shields.io/badge/GORM-ORM-336791)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/Auth-JWT-black)

## 🌍 Live URL

- Live API: Add your deployed URL here
- Local API: http://localhost:8080

## ✨ Core Features

- 🔐 JWT authentication with register, login, and current-user profile endpoints
- 🛡️ Role-based authorization for admin and driver permissions
- 🅿️ Parking zone management (public read, admin write)
- ⚡ Concurrency-safe reservation creation for limited-capacity zones
- 📦 Clean layered architecture with explicit DTO, handler, service, repository separation
- ✅ Request validation and consistent success/error response patterns

## 🛠️ Tech Stack

- Language: Go (module version in this project: 1.25.4)
- HTTP Framework: Echo v5
- ORM: GORM + PostgreSQL driver
- Database: PostgreSQL (local or hosted, e.g., Neon/Supabase)
- Auth: github.com/golang-jwt/jwt/v5
- Validation: github.com/go-playground/validator/v10
- Config: github.com/joho/godotenv

## 🏛️ Architecture (Layer Interaction)

The project follows strict separation of concerns:

1. DTO layer defines request and response contracts.
2. Handler layer handles HTTP, binding, validation, and response shaping.
3. Service layer holds business logic and permission-sensitive decisions.
4. Repository layer contains all database operations, transactions, and locks.
5. Models/entities represent database tables.

```mermaid
flowchart LR
	A[Client] --> B[Echo Router]
	B --> C[Middleware: Auth + Role]
	C --> D[Handler]
	D --> E[Service]
	E --> F[Repository]
	F --> G[(PostgreSQL)]
```

### Key Structure

- cmd/main.go
- internal/server/http.go
- internal/config
- internal/middlewares
- internal/domain/user
- internal/domain/zone
- internal/domain/reservation

## 🔒 Database Query Lock (Reservation Create)

To prevent overbooking under concurrent requests, reservation creation uses:

- A single DB transaction
- Row-level lock on the target parking zone using FOR UPDATE
- Atomic capacity check and insert inside the same transaction

Implemented in CreateWithCapacityLock (reservation repository):

- Lock zone row first
- Count active reservations in zone
- Reject if active count >= total capacity
- Reject duplicate active license plate inside the same zone
- Create reservation only when all checks pass

This solves the race condition where two clients try to reserve the final available spot at the same time.

## 🧪 Concurrency Testing (Create Reservation)

### k6 Test Script

- Script path: tmp/k6-reservations.js
- Generates dynamic license plates within validation limit
- Treats both 201 (created) and 409 (business conflict) as expected responses

Example run:

1. Set environment variables:
   - TOKEN
   - BASE_URL=http://127.0.0.1:8080
   - ZONE_ID=4
   - RATE=40
   - DURATION=30s
   - PRE_VUS=100
   - MAX_VUS=300
2. Run: k6 run ./tmp/k6-reservations.js

### hey Load Test

- Useful for quick stress checks
- For PowerShell, use JSON body file with -D to avoid quoting issues
- Prefer 127.0.0.1 over localhost to avoid IPv6 connection ambiguity

## ⚙️ Local Setup

### 1) Prerequisites

- Go 1.22+
- PostgreSQL

### 2) Install dependencies

- go mod tidy

### 3) Configure .env

Required variables:

| Variable | Required | Default | Notes |
|---|---|---|---|
| DSN | Yes | - | PostgreSQL connection string |
| PORT | No | 8080 | API listen port |
| JWT_SECRET | Yes | - | JWT signing secret |
| JWT_EXPIRY_HOURS | No | 24 | Must be > 0 |
| BCRYPT_COST | No | 10 | Must be 10-12 |

Example:

DSN="postgresql://user:password@host:5432/dbname?sslmode=disable"
PORT=8080
JWT_SECRET="secret_key"
JWT_EXPIRY_HOURS=24
BCRYPT_COST=10

### 4) Run server

- go run ./cmd/main.go
- Optional hot reload: air

### 5) Health check

- GET /health
- Expected: spot-sync running

## 🌐 API Endpoint List

Base path: /api/v1

### 🔹 Auth Module

| Method | Endpoint | Access |
|---|---|---|
| POST | /auth/register | Public |
| POST | /auth/login | Public |
| GET | /auth/me | Authenticated |

### 🔹 Parking Zones Module

| Method | Endpoint | Access |
|---|---|---|
| GET | /zones | Public |
| GET | /zones/:id | Public |
| POST | /zones | Admin |
| PUT | /zones/:id | Admin |
| DELETE | /zones/:id | Admin |

### 🔹 Reservations Module

| Method | Endpoint | Access |
|---|---|---|
| POST | /reservations | Authenticated (driver/admin) |
| GET | /reservations/my-reservations | Authenticated |
| DELETE | /reservations/:id | Authenticated |
| GET | /reservations | Admin |

## 🔐 Auth Header Format

Authorization: Bearer <token>

## 📝 Notes

- Auto-migration runs on startup for users, parking_zones, and reservations.
- Validation errors return 400.
- Unauthorized access returns 401; insufficient role returns 403.
- Reservation business conflicts (zone full, duplicate plate) return 409.

