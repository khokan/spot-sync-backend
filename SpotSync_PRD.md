# Product Requirements Document (PRD)
## SpotSync: Smart Parking and EV Charging Reservation Platform

## 1. Document Control
- Product Name: SpotSync
- Version: 1.0
- Date: 2026-06-26
- Prepared For: Engineering, QA, Product, DevOps, and Stakeholders
- Status: Final Draft for Build Execution

## 2. Executive Summary
SpotSync is a centralized platform for high-traffic venues such as airports and malls to manage parking zones and reservations, with a focus on scarce EV charging spots. The system enables drivers to discover availability and reserve spots while enabling administrators to manage zones, pricing, and reservation visibility.

The platform must enforce strict capacity constraints under concurrent demand and provide secure authentication and role-based access control. The solution is implemented in Go using Clean Architecture with Echo, GORM, PostgreSQL, JWT, bcrypt, and request validation.

## 3. Problem Statement
Venues with mixed parking inventory and limited EV charging capacity face:
- Reservation conflicts during peak demand.
- Overbooking risk caused by race conditions.
- Limited operational visibility for admins.
- Inconsistent user experience for drivers booking under time pressure.

SpotSync addresses these by delivering atomic reservation handling, transparent availability, and role-secured operations.

## 4. Product Vision and Objectives
### Vision
Provide a reliable, secure, and scalable reservation platform that prevents over-capacity booking while simplifying parking operations.

### Business Objectives
1. Reduce failed user booking attempts caused by stale availability.
2. Eliminate overbooking incidents in EV charging zones.
3. Improve administrative control over zone inventory and pricing.
4. Deliver auditable reservation lifecycle data for operations.

### Product Objectives
1. Support role-based workflows for driver and admin.
2. Guarantee capacity correctness under concurrent writes.
3. Expose clear, consistent API contracts for client apps.
4. Ensure production-grade security and observability.

## 5. Scope
### In Scope (MVP)
1. User registration and login with JWT.
2. Role-based authorization for driver and admin.
3. Parking zone CRUD (admin-controlled writes, public reads).
4. Reservation creation with strict transactional locking.
5. My reservations listing and cancellation.
6. Admin list of all reservations with user and zone data.
7. Standardized success and error response envelopes.

### Out of Scope (MVP)
1. Payment processing and invoicing.
2. Time-slot based future booking windows.
3. Push notifications, email, and SMS.
4. Dynamic surge pricing.
5. Multi-tenant organizations and hierarchy.
6. Mobile SDKs and web front-end implementation.

## 6. Stakeholders and Personas
### Stakeholders
1. Venue Operations Managers
2. Drivers (end users)
3. System Administrators
4. Engineering and QA teams

### Personas
1. Driver
- Needs quick visibility of available spots.
- Needs fast and reliable reservation flow.
- Needs control to view and cancel own reservations.

2. Admin
- Needs full control of parking zones and pricing.
- Needs global view of reservations for operations.
- Needs secure role-guarded access to privileged actions.

## 7. User Stories
1. As a driver, I can register and log in so I can access reservations.
2. As a driver, I can view all zones with real-time available spots.
3. As a driver, I can reserve a spot if capacity exists.
4. As a driver, I can view my reservation history and statuses.
5. As a driver, I can cancel my own reservation.
6. As an admin, I can create, update, and delete zones.
7. As an admin, I can set and adjust zone pricing.
8. As an admin, I can view all reservations with user and zone details.
9. As the system, I must reject over-capacity reservations under concurrency.

## 8. Functional Requirements

### FR-1: Authentication and Identity
1. The system shall allow public user registration with name, email, password, and role.
2. The system shall enforce unique email addresses.
3. The system shall hash passwords using bcrypt cost 10-12 before storage.
4. The system shall allow public login with email and password.
5. The system shall return JWT tokens upon successful login.
6. JWT payload shall include user id and role.
7. Passwords shall never be returned in API responses or logs.

### FR-2: Authorization
1. Protected endpoints shall require a valid JWT.
2. Requests without valid JWT shall return 401 Unauthorized.
3. Admin-only endpoints shall enforce role check and return 403 Forbidden if role is insufficient.
4. Role validation may occur in middleware or handler before service invocation.

### FR-3: Parking Zones
1. Admin shall create parking zones with required fields.
2. Admin shall update and delete zones.
3. Public users shall read all zones and single zone by id.
4. Zone list and details shall include available_spots computed dynamically.
5. available_spots shall be calculated as total_capacity minus active reservation count.

### FR-4: Reservations
1. Authenticated users (driver and admin) shall create reservations.
2. Reservation requires zone_id and license_plate (max 15 chars).
3. Reservation status defaults to active.
4. Users shall list their own reservations with zone details.
5. Users shall cancel reservations by id.
6. Driver can cancel only own reservations.
7. Admin can view all reservations with user and zone preloaded.

### FR-5: Concurrency and Capacity Integrity
1. Reservation creation shall run in a single database transaction.
2. Reservation creation shall lock target parking zone row using FOR UPDATE.
3. System shall count active reservations inside the transaction after lock acquisition.
4. If active count is equal to or greater than total_capacity, creation shall fail with 409 Conflict.
5. Under simultaneous competing requests, exactly one request may consume the last spot.

### FR-6: Validation and Error Handling
1. Request DTO validation shall use validator integration.
2. Invalid payloads shall return 400 Bad Request with validation details.
3. Unknown resources shall return 404 Not Found.
4. Duplicate conflicts or business conflicts shall return 409 Conflict.
5. Unexpected failures shall return 500 Internal Server Error.
6. All responses shall use standard success and error envelopes.

## 9. Non-Functional Requirements

### NFR-1: Architecture and Code Quality
1. Clean Architecture layering is mandatory.
2. Handlers must never directly access GORM or database connections.
3. DTOs must isolate API contracts from persistence models.
4. Dependency injection must be manual in main wiring:
- Repository to Service to Handler.

### NFR-2: Security
1. JWT signing secret shall be environment-driven and never hardcoded.
2. Password hashing shall use bcrypt with approved cost range.
3. Sensitive information shall be excluded from responses and logs.
4. API should enforce HTTPS in production environment.
5. Token expiration and validation must be strictly enforced.

### NFR-3: Performance
1. API should sustain peak read traffic for zone availability endpoints.
2. Reservation writes shall maintain correctness under high contention.
3. Query design should use indexes to support fast lookups.

### NFR-4: Reliability
1. Capacity invariants shall hold at all times.
2. Transaction rollback shall occur on reservation failure.
3. API shall remain consistent under partial failures.

### NFR-5: Observability
1. Structured logs shall include request id, endpoint, status code, and latency.
2. Errors shall be logged with correlation identifiers.
3. Metrics should include reservation attempt count, success/failure ratio, and conflict frequency.

### NFR-6: Maintainability
1. Public API contracts remain stable across patch releases.
2. Layer boundaries should enable isolated unit testing.
3. Business rules must be centralized in service layer.

## 10. Technical Architecture

### 10.1 Stack
1. Go 1.22+
2. Echo v4 for HTTP routing and middleware
3. GORM with PostgreSQL driver
4. PostgreSQL (NeonDB or Supabase compatible)
5. go-playground validator v10
6. golang-jwt jwt v5
7. x/crypto bcrypt

### 10.2 Layer Responsibilities
1. DTO layer
- Defines request and response structures.
- Applies validation tags.
- Prevents leakage of GORM models in API output.

2. Handler layer
- Parses and validates request DTOs.
- Extracts claims from Echo context.
- Performs authorization checks.
- Calls service methods and returns response envelopes.

3. Service layer
- Implements business rules.
- Hashes passwords and creates JWT.
- Enforces ownership and capacity checks.
- Orchestrates repository calls.

4. Repository layer
- Performs all data access via GORM.
- Encapsulates transactions and locking logic.
- Handles preloads and persistence concerns.

5. Models layer
- GORM entities and table mappings.

## 11. Data Model and Constraints

### users
1. id: auto-increment primary key
2. name: required
3. email: required, unique
4. password: required, hashed
5. role: enum-like constraint: driver or admin, default driver
6. created_at: auto-generated timestamp
7. updated_at: auto-updated timestamp

### parking_zones
1. id: auto-increment primary key
2. name: required
3. type: constrained to general, ev_charging, covered
4. total_capacity: integer, required, greater than 0
5. price_per_hour: decimal/float, required, greater than 0
6. created_at
7. updated_at

### reservations
1. id: auto-increment primary key
2. user_id: foreign key users.id, required
3. zone_id: foreign key parking_zones.id, required
4. license_plate: required, max length 15
5. status: constrained to active, completed, cancelled, default active
6. created_at
7. updated_at

### Suggested Indexes
1. users.email unique index
2. reservations.zone_id index
3. reservations.user_id index
4. reservations.status index
5. composite index on reservations(zone_id, status) for capacity checks

## 12. API Contract Requirements

### 12.1 Base Path
- /api/v1

### 12.2 Authentication Endpoints
1. POST /auth/register
- Public
- Creates account
- Returns created user profile (without password)

2. POST /auth/login
- Public
- Returns JWT token plus user summary

### 12.3 Parking Zone Endpoints
1. POST /zones
- Admin only
- Creates zone

2. GET /zones
- Public
- Lists zones with available_spots

3. GET /zones/:id
- Public
- Returns one zone with available_spots

### 12.4 Reservation Endpoints
1. POST /reservations
- Authenticated
- Creates reservation with transactional lock

2. GET /reservations/my-reservations
- Authenticated
- Returns caller-owned reservations with zone preload

3. DELETE /reservations/:id
- Authenticated
- Cancels reservation
- Driver limited to own reservation
- Admin allowed per policy decision (recommended yes)

4. GET /reservations
- Admin only
- Lists all reservations with user and zone preload

### 12.5 Response Envelope Standard
1. Success:
- success: true
- message: operation description
- data: payload

2. Error:
- success: false
- message: error description
- errors: details

### 12.6 HTTP Status Mapping
1. 200 for successful reads and deletes
2. 201 for created resources
3. 400 for invalid input or validation failure
4. 401 for authentication failure
5. 403 for authorization failure
6. 404 for missing resources
7. 409 for business conflicts such as zone full
8. 500 for unhandled server/database errors

## 13. Core Business Rules
1. Email is globally unique.
2. Password must be hashed before persistence.
3. Role defaults to driver if omitted.
4. Zone capacity must always remain non-negative in availability calculations.
5. Only active reservations consume capacity.
6. Canceling reservation changes status to cancelled, freeing capacity.
7. Drivers cannot mutate resources they do not own where ownership applies.
8. Zone full condition must return deterministic conflict response.

## 14. Concurrency Control Specification (Critical)
1. Reservation creation must execute as an atomic transaction.
2. Target zone row must be selected with UPDATE lock.
3. Active reservation count must be computed inside same transaction.
4. If capacity available, create reservation and commit.
5. If capacity exceeded, return zone full error and rollback.
6. Repository must expose a dedicated transactional method for this flow.
7. Service must map repository capacity conflict to 409 response contract.

Acceptance requirement:
- In a load test with N parallel requests competing for one remaining spot, exactly one succeeds and N-1 fail with 409.

## 15. Security and Compliance Requirements
1. JWT secret, DB credentials, and bcrypt cost must be configurable via environment variables.
2. Token includes minimal claims: user id, role, expiry.
3. Middleware validates token signature and expiration.
4. No plaintext password storage under any condition.
5. Logs must avoid secrets and sensitive payload fields.
6. CORS policy must be explicitly configured for approved clients.

## 16. Logging, Monitoring, and Alerting
1. Log all requests with method, path, status, latency, and request id.
2. Log authentication failures and authorization denials.
3. Track metrics:
- reservation_create_total
- reservation_create_conflict_total
- zone_available_spots_query_latency
- auth_login_failure_total
4. Alert on sudden spikes in 500 or 409 rates beyond baseline.
5. Add health endpoint for readiness/liveness checks.

## 17. Test Strategy and Acceptance

### 17.1 Unit Tests
1. Service tests for registration, login, role checks, and ownership checks.
2. Reservation service tests for zone full and happy path.
3. DTO validation tests for required fields and constraints.

### 17.2 Integration Tests
1. End-to-end auth flow with JWT middleware.
2. Zone CRUD permission matrix.
3. Reservation creation and cancellation lifecycle.
4. My reservations preload behavior.

### 17.3 Concurrency Tests
1. Simulate simultaneous reservation requests against last spot.
2. Validate no over-capacity outcome.
3. Confirm exact conflict response behavior.

### 17.4 Exit Criteria
1. All API endpoints pass contract tests.
2. Concurrency test proves no overbooking.
3. Security checks pass for password and token handling.
4. Error envelope and status code matrix are consistent.

## 18. Delivery Plan

### Phase 1: Foundation
1. Project structure with strict layer boundaries.
2. Database models and migrations.
3. DI wiring and middleware setup.

### Phase 2: Auth and Zones
1. Registration and login.
2. JWT middleware and role middleware.
3. Zone create/update/delete/list/detail.

### Phase 3: Reservation Core
1. Transactional lock-based reservation create.
2. My reservations and cancel flow.
3. Admin reservations listing.

### Phase 4: Hardening
1. Automated tests including concurrency scenario.
2. Observability and error handling improvements.
3. Deployment readiness checklist.

## 19. Risks and Mitigations
1. Risk: Overbooking under race condition.
- Mitigation: Mandatory transaction and row lock.

2. Risk: Role bypass bugs.
- Mitigation: Centralized auth middleware and handler-level guard checks.

3. Risk: Performance degradation on zone listing.
- Mitigation: Proper indexing and optimized availability queries.

4. Risk: Inconsistent API responses.
- Mitigation: Shared response builder and contract tests.

5. Risk: Secret leakage.
- Mitigation: Environment-based config and sanitized logging.

## 20. Open Decisions
1. Should admins be allowed to cancel any reservation through existing cancel endpoint or a separate admin endpoint.
2. Whether completed status transitions are manual or future automated workflow.
3. Final JWT expiration policy and refresh token strategy for future versions.

## 21. Definition of Done
1. All listed endpoints implemented per contract.
2. Clean Architecture constraints fully respected.
3. Concurrency rule implemented and proven by test.
4. Security requirements validated.
5. Documentation includes setup, environment variables, and API examples.
6. Production deployment checklist completed.
