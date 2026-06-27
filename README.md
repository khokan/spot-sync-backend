# SpotSync Backend

SpotSync Backend is a Go-based REST API for managing parking zones and vehicle reservations with role-based authorization.

## Live URL

- Live API URL: Not deployed yet (add your production/staging URL here)
- Local API URL: `http://localhost:8080`

## Features

- JWT-based authentication (`register`, `login`, `me`)
- Role-based access control (`admin`, `driver`)
- Parking zone management (admin CRUD + public read)
- Reservation lifecycle (create, list own, cancel, admin list all)
- Input validation using request DTOs and validator rules
- Layered architecture (handler, service, repository, entity/dto)

## Tech Stack

- Language: Go `1.25.4`
- HTTP framework: Echo v5
- ORM: GORM
- Database: PostgreSQL
- Auth: JWT (`github.com/golang-jwt/jwt/v5`)
- Validation: `go-playground/validator/v10`
- Config loading: `joho/godotenv`

## Architecture

This project follows a layered architecture to keep transport logic, business rules, and data access separate.

### Layer Interaction

1. **Route/Register Layer** wires dependencies and route groups.
2. **Handler Layer** accepts HTTP requests, binds/validates payload, and returns responses.
3. **Service Layer** contains business rules and orchestration.
4. **Repository Layer** performs DB operations through GORM.
5. **Entity/DTO Layer** defines persistence models and request/response contracts.

```mermaid
flowchart TD
	A[Client] --> B[Echo Router]
	B --> C[Middleware<br/>Auth + Role]
	C --> D[Handler]
	D --> E[Service]
	E --> F[Repository]
	F --> G[(PostgreSQL)]

	D --> H[DTO Validation]
	E --> I[Business Rules]
```

### Project Structure (Key Paths)

```text
cmd/main.go                     # Application entrypoint
internal/server/http.go         # Echo server bootstrap + route registration
internal/config/                # Env config + DB connection
internal/auth/                  # JWT service
internal/middlewares/           # Auth and role middleware
internal/domain/user/           # User module (auth)
internal/domain/zone/           # Zone module
internal/domain/reservation/    # Reservation module
```

## Setup (Run Locally)

### 1. Prerequisites

- Go 1.25+
- PostgreSQL database

### 2. Clone and install dependencies

```bash
git clone <your-repo-url>
cd spot-sync-backend
go mod tidy
```

### 3. Configure environment variables

Create/update `.env` in project root:

```env
DSN="postgresql://user:password@host:5432/dbname?sslmode=disable"
PORT=8080
BCRYPT_COST=10
JWT_SECRET="secret_key"
JWT_EXPIRY_HOURS=24
```

### 4. Run the server

```bash
go run ./cmd/main.go
```

Optional hot reload (if `air` is installed):

```bash
air
```

### 5. Health check

```http
GET /health
```

Response:

```text
spot-sync running
```

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DSN` | Yes | - | PostgreSQL connection string |
| `PORT` | No | `8080` | API server port |
| `JWT_SECRET` | Yes | - | Secret key for signing JWT |
| `JWT_EXPIRY_HOURS` | No | `24` | JWT expiry duration in hours |
| `BCRYPT_COST` | No | `10` | Password hashing cost (`10` to `12`) |

## API Endpoints

Base prefix: `/api/v1`

### Auth (`/api/v1/auth`)

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/auth/register` | Public | Register a new user |
| `POST` | `/auth/login` | Public | Authenticate and receive JWT |
| `GET` | `/auth/me` | Authenticated | Get current user profile |

### Zones (`/api/v1/zones`)

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/zones` | Public | List all parking zones |
| `GET` | `/zones/:id` | Public | Get zone details by ID |
| `POST` | `/zones` | Admin | Create a parking zone |
| `PUT` | `/zones/:id` | Admin | Update a parking zone |
| `DELETE` | `/zones/:id` | Admin | Delete a parking zone |

### Reservations (`/api/v1/reservations`)

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/reservations` | Authenticated | Create reservation |
| `GET` | `/reservations/my-reservations` | Authenticated | List current user's reservations |
| `DELETE` | `/reservations/:id` | Authenticated | Cancel reservation (owner/admin rules apply) |
| `GET` | `/reservations` | Admin | List all reservations |

## Authentication

Send JWT in `Authorization` header:

```http
Authorization: Bearer <token>
```

## Notes

- Auto migration runs on startup for `User`, `ParkingZone`, and `Reservation`.
- Validation errors return HTTP `400`.
- Protected endpoints return HTTP `401` for invalid/missing token and `403` for insufficient role.

