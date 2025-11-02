---
name: Bookings API - Phase 1
about: Project Setup for bookings-api
title: 'Phase 1: Project Setup - bookings-api'
labels: enhancement, bookings-api, phase-1
assignees: ''
---

## Phase 1: Project Setup

**Service:** bookings-api
**Priority:** High
**Estimated Time:** 2-3 hours
**Branch:** `feature/bookings-api-phase-1-setup`

---

### Description
Initialize the bookings-api project with all necessary dependencies, configuration files, and basic server structure. This follows the same pattern as users-api.

---

### Tasks
- [ ] Initialize Go modules (`go mod init bookings-api`)
- [ ] Create `.env` file with all required variables
- [ ] Create `.env.example` as template
- [ ] Create `cmd/api/main.go` with basic HTTP server
- [ ] Add all dependencies to go.mod:
  - github.com/gin-gonic/gin v1.10.0
  - github.com/golang-jwt/jwt/v5 v5.2.1
  - github.com/joho/godotenv v1.5.1
  - github.com/google/uuid v1.6.0
  - github.com/rs/zerolog v1.33.0
  - github.com/streadway/amqp v1.1.0
  - gorm.io/gorm v1.25.12
  - gorm.io/driver/mysql v1.5.7
- [ ] Verify compilation with `go build`

---

### Success Criteria
- [ ] `go mod tidy` runs successfully
- [ ] `go build ./cmd/api` compiles without errors
- [ ] Server starts on port 8003
- [ ] Basic health endpoint responds with 200 OK
- [ ] All dependencies downloaded correctly

---

### Dependencies
**Requires:**
- MySQL running (can use users-api MySQL instance)
- RabbitMQ installed/running

**Blocks:**
- Phase 2 (Database Setup)

---

### Files to Create
```
backend/bookings-api/
├── cmd/api/main.go
├── .env
├── .env.example
├── go.mod
└── go.sum (auto-generated)
```

---

### .env Template
```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=Prueba.9876
DB_NAME=carpooling_bookings

# JWT (same secret as users-api)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server
SERVER_PORT=8003

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin@localhost:5672/

# External APIs
TRIPS_API_URL=http://localhost:8002
USERS_API_URL=http://localhost:8001
```

---

### Testing Steps
```bash
# Navigate to project
cd backend/bookings-api

# Verify dependencies
go mod tidy

# Compile
go build ./cmd/api

# Run server
go run cmd/api/main.go

# In another terminal, test health endpoint
curl http://localhost:8003/health
# Expected: {"status":"ok"} or similar
```

---

### Implementation Guide
1. Create branch: `git checkout -b feature/bookings-api-phase-1-setup`
2. Use Claude Code plan mode:
   ```
   @CONTEXT_BOOKINGS_API.md Implement Phase 1 (Project Setup) for bookings-api.

   Tasks:
   - Initialize go.mod
   - Setup .env file
   - Create main.go with basic server
   - Add all dependencies

   Follow pattern from backend/users-api/cmd/api/main.go
   ```
3. Review and approve plan
4. Verify all success criteria
5. Commit: `git commit -m "feat(bookings): setup project structure and dependencies"`
6. Push: `git push -u origin feature/bookings-api-phase-1-setup`
7. Create PR to `dev`

---

### References
- Implementation Plan: `GITFLOW.md` - Phase 1
- Context: `CONTEXT_BOOKINGS_API.md`
- Pattern Reference: `backend/users-api/cmd/api/main.go`
- How to Use: `HOW_TO_USE_PLAN_MODE.md`

---

### Notes
- Use same JWT_SECRET as users-api for consistency
- DB_NAME should be `carpooling_bookings` (not `carpooling_users`)
- Port 8003 (not 8001 like users-api)
- Make sure MySQL is running before testing
