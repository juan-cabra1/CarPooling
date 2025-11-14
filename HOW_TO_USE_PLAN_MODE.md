# How to Use Claude Code Plan Mode

## Quick Start

### 1. Navigate to Project
```bash
cd /home/user/CarPooling
```

### 2. Create Feature Branch
```bash
# Always start from latest dev
git checkout dev
git pull origin dev

# Create feature branch (example for trips-api issue #1)
git checkout -b feature/trips-api/1-project-setup
```

### 3. Start Claude Code
```bash
claude
```

### 4. Use Plan Mode with Context

```
@CONTEXT_TRIPS_API.md Implement Phase 1 (Project Setup) for trips-api.

Tasks:
- Initialize go mod init trips-api
- Create .env with MongoDB, RabbitMQ, JWT config
- Setup cmd/api/main.go with basic server on port 8002
- Add dependencies: gin, mongo-driver, uuid, zerolog, amqp

Follow pattern from backend/users-api/cmd/api/main.go but for MongoDB
```

---

## Pattern for Each Feature

```
@CONTEXT_{SERVICE}_API.md Implement [Feature Name] for {service}-api.

Tasks:
- Task 1
- Task 2
- Task 3

[Any additional context or patterns to follow]
```

---

## Example Prompts

### trips-api Phase 1: Project Setup
```
@CONTEXT_TRIPS_API.md Implement Phase 1 (Project Setup) for trips-api.

Initialize the project with:
- Go modules (go mod init trips-api)
- .env file with MONGO_URI, JWT_SECRET, SERVER_PORT, RABBITMQ_URL
- cmd/api/main.go with Gin server on port 8002
- All dependencies from CONTEXT_TRIPS_API.md

Reference backend/users-api structure
```

### trips-api Phase 2: MongoDB Connection
```
@CONTEXT_TRIPS_API.md Implement MongoDB connection and domain models for trips-api.

Create:
1. internal/config/config.go - Load env variables
2. internal/domain/trip.go - Trip struct with all fields from MongoDB schema
3. Connect to MongoDB in main.go
4. Create indexes: driver_id, status, departure_datetime, origin.city+destination.city

MongoDB schema is in CONTEXT_TRIPS_API.md
```

### trips-api Phase 3: Repository Layer
```
@CONTEXT_TRIPS_API.md Implement repository layer with MongoDB driver.

Create:
1. internal/repository/trip_repository.go
   - Interface with CRUD methods
   - MongoDB implementation with mongo-driver

2. internal/repository/event_repository.go
   - IsEventProcessed(eventID) (bool, error)
   - MarkEventProcessed(event) error
   - Use UNIQUE index on event_id for idempotency

Include optimistic locking in UpdateAvailability method
```

### trips-api Phase 7: RabbitMQ Consumer (CRITICAL)
```
@CONTEXT_TRIPS_API.md Implement RabbitMQ consumer for reservation events - CRITICAL PHASE.

This implements idempotency to prevent duplicate seat decreases.

Create internal/messaging/reservation_consumer.go with:
- HandleReservationEvent(event) error
- IDEMPOTENCY check using CheckAndMarkEvent() BEFORE processing
- handleReservationCreated(event) - decrease seats with optimistic lock
- handleReservationCancelled(event) - increase seats

Critical:
- Check idempotency FIRST (event_id in processed_events collection)
- Skip if already processed
- Manual ACK only after success
- NACK with requeue on transient errors
- Use optimistic locking (availability_version field)

Event format and idempotency logic in CONTEXT_TRIPS_API.md
```

---

## What Happens in Plan Mode

1. **Claude reads context** - Loads CONTEXT_TRIPS_API.md
2. **Analyzes patterns** - Reviews users-api for reference
3. **Creates plan** - Detailed step-by-step implementation
4. **Asks for approval** - Reviews plan with you
5. **Implements** - Executes the plan after approval

---

## During Plan Mode

### Good Practices
✅ Be specific about which feature/phase
✅ Reference context file with @
✅ Mention patterns to copy (e.g., from users-api)
✅ Include success criteria from issue
✅ Ask questions if plan is unclear

### What to Avoid
❌ Generic prompts ("implement trips-api")
❌ Forgetting to reference context file
❌ Approving plan without understanding it
❌ Skipping testing steps

---

## After Implementation

### 1. Test
```bash
cd backend/trips-api

# Compile
go mod tidy
go build ./cmd/api

# Run tests (if applicable)
go test ./... -v

# Manual test
go run cmd/api/main.go
# In another terminal:
curl http://localhost:8002/health
```

### 2. Commit
```bash
git status
git diff

git add .
git commit -m "feat(trips): initialize project structure and dependencies

- Create go.mod and go.sum
- Setup .env configuration
- Add basic main.go with Gin server
- Install all dependencies

Closes #1"
```

### 3. Push and PR
```bash
git push -u origin feature/trips-api/1-project-setup

# Create PR on GitHub:
# Base: dev
# Compare: feature/trips-api/1-project-setup
# Title: feat(trips): Project Setup - Issue #1
```

---

## Workflow Summary

```
1. git checkout dev && git pull
2. git checkout -b feature/trips-api/{issue}-{name}
3. claude
4. @CONTEXT_TRIPS_API.md Implement {feature}...
5. [Review plan, approve]
6. [Claude implements]
7. Test (compile, run, curl)
8. git commit -m "feat(trips): {description}"
9. git push -u origin feature/...
10. Create PR to dev
11. After merge, repeat from step 1
```

---

## Common Issues

### "Claude can't find context file"
```bash
# Use absolute path
@/home/user/CarPooling/CONTEXT_TRIPS_API.md
```

### "Plan is too vague"
```bash
# Be more specific
Instead of: "Implement trips"
Use: "Implement repository layer with MongoDB driver following CONTEXT_TRIPS_API.md schema"
```

### "Tests failing"
```bash
# Run specific test with verbose
go test ./internal/service/idempotency_service_test.go -v

# Check MongoDB connection
docker ps | grep mongo
```

---

## Tips

1. **One feature at a time** - Don't mix multiple phases
2. **Read context first** - Understand what you're building
3. **Review plan carefully** - Especially for critical features (idempotency, optimistic locking)
4. **Test before committing** - Verify compilation and basic functionality
5. **Small commits** - One logical change per commit
6. **Reference issues** - Use "Closes #N" in commit message

---

## Quick Reference

| Command | Purpose |
|---------|---------|
| `claude` | Start Claude Code |
| `@CONTEXT_TRIPS_API.md` | Reference context in prompt |
| `@backend/users-api` | Reference existing code |
| `/exit` | Exit Claude Code |
| `git checkout -b feature/trips-api/N-name` | Create branch |
| `go mod tidy` | Clean dependencies |
| `go build ./cmd/api` | Compile |
| `go test ./... -v` | Run tests |
| `git push -u origin feature/...` | Push branch |

---

Happy coding! 🚀

For more details, see:
- `CONTEXT_TRIPS_API.md` - Complete specification
- `GITFLOW.md` - Git workflow
- `README_DEVELOPMENT.md` - Project overview
