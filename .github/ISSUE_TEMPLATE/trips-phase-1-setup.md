---
name: trips-api - Phase 1 Setup
about: Initialize trips-api project with MongoDB
title: '[trips-api] Project Setup'
labels: trips-api, enhancement, phase-1
assignees: ''
---

## Phase 1: Project Setup

**Service:** trips-api
**Priority:** High
**Estimated Time:** 2-3 hours
**Branch:** `feature/trips-api/1-project-setup`

---

### Description
Initialize the trips-api project with Go modules, MongoDB connection, and basic server structure. This is the main API of the platform (faculty requirement).

---

### Tasks
- [ ] Initialize Go modules (`go mod init trips-api`)
- [ ] Create `.env` file with MongoDB, RabbitMQ, JWT config
- [ ] Create `.env.example` as template
- [ ] Create `cmd/api/main.go` with basic HTTP server (port 8002)
- [ ] Add dependencies:
  - github.com/gin-gonic/gin
  - go.mongodb.org/mongo-driver
  - github.com/golang-jwt/jwt/v5
  - github.com/joho/godotenv
  - github.com/google/uuid
  - github.com/rs/zerolog
  - github.com/streadway/amqp
- [ ] Verify compilation with `go build`

---

### Success Criteria
- [ ] `go mod tidy` runs successfully
- [ ] `go build ./cmd/api` compiles without errors
- [ ] Server starts on port 8002
- [ ] Basic health endpoint responds with 200 OK
- [ ] All dependencies downloaded correctly
- [ ] MongoDB connection string configured

---

### Dependencies
**Requires:**
- MongoDB running (use docker-compose from search-api)

**Blocks:**
- Phase 2 (MongoDB Connection & Models)

---

### Files to Create
```
backend/trips-api/
├── cmd/api/main.go
├── .env
├── .env.example
├── go.mod
└── go.sum (auto-generated)
```

---

### .env Template
```bash
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=carpooling_trips

# JWT (same secret as users-api)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server
SERVER_PORT=8002

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin@localhost:5672/

# External APIs
USERS_API_URL=http://localhost:8001
```

---

### Testing Steps
```bash
# Navigate to project
cd backend/trips-api

# Verify dependencies
go mod tidy

# Compile
go build ./cmd/api

# Run server
go run cmd/api/main.go

# In another terminal, test health endpoint
curl http://localhost:8002/health
# Expected: {"status":"ok"} or similar

# Verify MongoDB connection (if implemented in Phase 1)
# Check logs for "MongoDB connected" message
```

---

### Implementation Guide
1. Create branch: `git checkout -b feature/trips-api/1-project-setup`
2. Use Claude Code plan mode:
   ```
   @CONTEXT_TRIPS_API.md Implement Phase 1 (Project Setup) for trips-api.

   Tasks:
   - Initialize go.mod for trips-api
   - Setup .env file with MongoDB config
   - Create main.go with basic Gin server on port 8002
   - Add all dependencies

   Reference backend/users-api/cmd/api/main.go for patterns
   ```
3. Review and approve plan
4. Verify all success criteria
5. Commit: `git commit -m "feat(trips): initialize project structure and dependencies"`
6. Push: `git push -u origin feature/trips-api/1-project-setup`
7. Create PR to `dev`

---

### References
- Context: `CONTEXT_TRIPS_API.md`
- Pattern Reference: `backend/users-api/cmd/api/main.go`
- Workflow: `GITFLOW.md`

---

### Notes
- Use same JWT_SECRET as users-api for consistency
- Port 8002 (not 8001 like users-api)
- MongoDB URI points to `carpooling_trips` database
- Make sure MongoDB is running before testing (use docker-compose from search-api)
