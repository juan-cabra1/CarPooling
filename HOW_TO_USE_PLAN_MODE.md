# How to Use Plan Mode with Claude Code

## Quick Start Guide

### 1. Initial Setup (One Time)

```bash
cd /home/user/CarPooling

# Verify context files exist
ls -la CONTEXT_BOOKINGS_API.md GITFLOW.md

# Checkout dev branch
git checkout dev
git pull origin dev
```

---

## 2. Start Working on a Phase

### Step 1: Create Feature Branch
```bash
# Example: Starting Phase 1
git checkout -b feature/bookings-api-phase-1-setup
```

### Step 2: Enter Plan Mode with Context

**Option A: Using file reference (RECOMMENDED)**
```bash
# From project root
claude

# In Claude Code, type:
@CONTEXT_BOOKINGS_API.md @GITFLOW.md I want to implement Phase 1 of bookings-api. Please create a plan for:
- Initializing go.mod
- Setting up .env file
- Creating main.go skeleton
- Adding all dependencies

Follow the structure from users-api at backend/users-api
```

**Option B: Inline context (if file refs don't work)**
```bash
claude

# Then paste the relevant section from CONTEXT_BOOKINGS_API.md
# For Phase 1, copy the "Phase 1 - Project Setup" section from GITFLOW.md
# Then say:

Based on this context, create a plan to implement Phase 1 of bookings-api
```

**Option C: Direct prompt with full context**
```bash
claude

# Type or paste:
I'm implementing bookings-api for a CarPooling project.

Context:
- Technology: Go 1.21, Gin, GORM, MySQL
- Port: 8003
- Pattern: Copy from backend/users-api
- Dependencies: gin, gorm, uuid, zerolog, amqp

Phase 1 Tasks:
1. Initialize go mod init bookings-api
2. Create .env with DB credentials, JWT secret, RabbitMQ URL
3. Create cmd/api/main.go with basic server
4. Add all dependencies

Please create a step-by-step plan
```

---

## 3. Example Prompts for Each Phase

### Phase 1: Project Setup
```
@CONTEXT_BOOKINGS_API.md I need to implement Phase 1 (Project Setup) for bookings-api.

Tasks:
- Initialize Go modules
- Create .env file with MySQL, RabbitMQ, JWT config
- Setup main.go skeleton
- Add dependencies: gin, gorm, uuid, zerolog, amqp

Follow the same structure as backend/users-api
```

### Phase 2: Database Setup
```
@CONTEXT_BOOKINGS_API.md I need to implement Phase 2 (Database Setup) for bookings-api.

Create:
1. internal/config/config.go (copy pattern from users-api)
2. internal/dao/booking.go with GORM tags
3. internal/dao/processed_event.go with GORM tags
4. internal/domain/booking.go with DTOs

Database schema is in CONTEXT_BOOKINGS_API.md

Ensure:
- UNIQUE index on processed_events.event_id
- All GORM tags match MySQL schema
- Auto-migration in main.go
```

### Phase 3: Repository Layer
```
@CONTEXT_BOOKINGS_API.md Implement Phase 3 (Repository Layer) for bookings-api.

Create:
1. internal/repository/booking_repository.go
   - Interface with CRUD methods
   - GORM implementation

2. internal/repository/event_repository.go
   - IsEventProcessed(eventID) (bool, error)
   - MarkEventProcessed(event) error
   - Use UNIQUE constraint for idempotency

Copy pattern from users-api repository layer
```

### Phase 7: RabbitMQ Consumer (CRITICAL)
```
@CONTEXT_BOOKINGS_API.md Implement Phase 7 (RabbitMQ Consumer) - CRITICAL PHASE

This is the most important phase. Implement:

1. internal/messaging/rabbitmq.go - connection setup
2. internal/messaging/trip_consumer.go with:
   - HandleTripEvent(event) error
   - IDEMPOTENCY check using CheckAndMarkEvent()
   - handleTripUpdated(event) error
   - handleTripCancelled(event) error
   - handleReservationFailed(event) error

Critical requirements:
- Check idempotency BEFORE processing
- Skip if already processed
- Manual ACK only after success
- NACK with requeue on transient errors
- Structured logging with zerolog

Event format in CONTEXT_BOOKINGS_API.md
```

### Phase 8: RabbitMQ Publisher
```
@CONTEXT_BOOKINGS_API.md Implement Phase 8 (RabbitMQ Publisher)

Create internal/messaging/reservation_publisher.go with:
- PublishReservationCreated(booking) error
- PublishReservationCancelled(booking, reason) error
- PublishReservationFailed(booking, reason) error

CRITICAL:
- Generate UUID for event_id BEFORE publishing (uuid.New().String())
- Set DeliveryMode: Persistent
- Set MessageId: event.EventID
- Publish to "reservations.events" exchange

Event format in CONTEXT_BOOKINGS_API.md
```

---

## 4. During Plan Mode

### What Claude Code Will Do:
1. Read the context files
2. Analyze users-api structure for patterns
3. Create a detailed step-by-step plan
4. Ask for your approval

### You Should:
1. Review the plan carefully
2. Ask questions if something is unclear
3. Request changes if needed
4. Approve to exit plan mode and start implementation

---

## 5. After Plan Approval (Implementation)

Claude will implement the plan. Monitor:

```bash
# Watch files being created
watch -n 1 'tree backend/bookings-api'

# Check git status frequently
git status

# View changes
git diff
```

---

## 6. Testing Each Phase

### Compile Check
```bash
cd backend/bookings-api
go mod tidy
go build ./cmd/api
```

### Run Tests
```bash
go test ./internal/service/... -v
go test ./... -cover
```

### Manual Testing
```bash
# Start server
go run cmd/api/main.go

# In another terminal, test endpoints
curl http://localhost:8003/health

# Test creating booking (after Phase 12)
curl -X POST http://localhost:8003/bookings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"trip_id":"507f1f77bcf86cd799439011","seats_reserved":2}'
```

---

## 7. Completing a Phase

### Commit Changes
```bash
cd /home/user/CarPooling

# Review changes
git status
git diff

# Stage all changes
git add backend/bookings-api

# Commit with conventional commit message
git commit -m "feat(bookings): implement phase 1 - project setup

- Initialize Go modules
- Setup .env configuration
- Create main.go skeleton
- Add all dependencies (gin, gorm, uuid, zerolog, amqp)

Closes #1"

# Push to remote
git push -u origin feature/bookings-api-phase-1-setup
```

### Create Pull Request
```bash
# Option 1: Using GitHub web interface
# Go to: https://github.com/juan-cabra1/CarPooling
# Click "Compare & pull request"
# Base: dev, Compare: feature/bookings-api-phase-1-setup
# Title: "Phase 1: Project Setup for bookings-api"
# Description: Copy success criteria from GITFLOW.md

# Option 2: Using gh CLI (if available)
gh pr create \
  --base dev \
  --head feature/bookings-api-phase-1-setup \
  --title "Phase 1: Project Setup for bookings-api" \
  --body "$(cat <<EOF
## Phase 1: Project Setup

### Completed Tasks
- [x] Initialize go.mod
- [x] Create .env and .env.example
- [x] Setup main.go skeleton
- [x] Add all dependencies

### Success Criteria
- [x] go mod tidy succeeds
- [x] go build succeeds
- [x] Server starts on port 8003

### Testing
\`\`\`bash
cd backend/bookings-api
go mod tidy
go build ./cmd/api
go run cmd/api/main.go
\`\`\`

Closes #1
EOF
)"
```

---

## 8. Moving to Next Phase

After PR is merged:

```bash
# Update dev branch
git checkout dev
git pull origin dev

# Start next phase
git checkout -b feature/bookings-api-phase-2-database

# Enter plan mode again
claude

# Use context for Phase 2
@CONTEXT_BOOKINGS_API.md Implement Phase 2 (Database Setup) for bookings-api...
```

---

## Tips for Effective Plan Mode Usage

### DO:
✅ Reference context files with @filename
✅ Be specific about which phase you're working on
✅ Mention patterns to copy from users-api
✅ Include success criteria in prompts
✅ Ask Claude to explain if something is unclear
✅ Review the plan before approving

### DON'T:
❌ Start coding without a plan
❌ Skip phases (they build on each other)
❌ Forget to reference CONTEXT_BOOKINGS_API.md
❌ Mix multiple phases in one branch
❌ Approve plan without understanding it
❌ Skip testing after implementation

---

## Example Full Workflow

```bash
# Phase 1 Start
git checkout dev
git pull origin dev
git checkout -b feature/bookings-api-phase-1-setup

# Enter Claude Code
claude

# In Claude:
# """
# @CONTEXT_BOOKINGS_API.md @GITFLOW.md
#
# Implement Phase 1 (Project Setup) for bookings-api.
# Follow patterns from backend/users-api
#
# Tasks from GITFLOW.md Issue #1:
# - Initialize go.mod
# - Create .env file
# - Setup main.go
# - Add dependencies
# """

# [Claude creates plan, you review and approve]
# [Claude implements]

# Verify
cd backend/bookings-api
go mod tidy
go build ./cmd/api
go run cmd/api/main.go
# (Ctrl+C to stop)

# Commit
git add .
git commit -m "feat(bookings): implement phase 1 - project setup"
git push -u origin feature/bookings-api-phase-1-setup

# Create PR (GitHub web or gh CLI)
# Wait for review/approval
# Merge to dev

# Phase 2 Start
git checkout dev
git pull origin dev
git checkout -b feature/bookings-api-phase-2-database
# ... repeat
```

---

## Troubleshooting

### "Claude can't find context file"
```bash
# Verify file exists
ls -la CONTEXT_BOOKINGS_API.md

# Try absolute path
@/home/user/CarPooling/CONTEXT_BOOKINGS_API.md

# Or paste content directly
```

### "Plan is too vague"
```
# Be more specific
Instead of: "Implement Phase 2"
Use: "Implement Phase 2 with BookingDAO GORM model using tags from CONTEXT_BOOKINGS_API.md schema"
```

### "Implementation differs from users-api pattern"
```
# Explicitly reference
"Copy the exact pattern from backend/users-api/internal/repository/user.go"
```

### "Tests failing after implementation"
```bash
# Run specific test
go test ./internal/service/idempotency_service_test.go -v

# Check logs
cat /tmp/bookings-api.log

# Ask Claude to fix
"""
Tests are failing with error: [paste error]
Please fix the issue in [file]
"""
```

---

## Advanced: Multi-Phase Planning

If you want to plan multiple phases at once:

```
@CONTEXT_BOOKINGS_API.md

I want to plan Phases 1-3 together (Setup, Database, Repository).

Create a comprehensive plan that:
1. Initializes the project (Phase 1)
2. Sets up database models and connection (Phase 2)
3. Implements repository layer (Phase 3)

Show dependencies between phases and order of implementation.
I'll implement in separate branches but want to understand the full picture.
```

---

## Context File Updates

If you need to update context during development:

```bash
# Edit context file
nano CONTEXT_BOOKINGS_API.md

# Commit changes
git add CONTEXT_BOOKINGS_API.md
git commit -m "docs: update bookings-api context with new requirements"

# Use updated context in next phase
claude
@CONTEXT_BOOKINGS_API.md [your prompt]
```

---

## Quick Reference Card

| Command | Purpose |
|---------|---------|
| `claude` | Start Claude Code |
| `@CONTEXT_BOOKINGS_API.md` | Reference context file |
| `@backend/users-api` | Reference existing code |
| `/exit` | Exit Claude Code |
| `/help` | Claude Code help |
| `git checkout -b feature/bookings-api-phase-N-name` | Start phase branch |
| `go mod tidy` | Clean dependencies |
| `go build ./cmd/api` | Verify compilation |
| `go test ./... -v` | Run tests |
| `git push -u origin feature/...` | Push branch |

---

## Success Indicators

You're doing it right if:
✅ Each phase takes 2-8 hours
✅ You understand the plan before approving
✅ Code compiles after each phase
✅ Tests pass after implementation
✅ PRs are small and focused
✅ You can explain what each file does
✅ Git history is clean with good commit messages

You need to adjust if:
❌ Phase takes more than 1 day
❌ You approve plan without understanding
❌ Code doesn't compile
❌ Tests fail and you don't know why
❌ PRs have 50+ files changed
❌ You can't explain the code
❌ Commits are "wip" or "fix"

---

Happy coding! 🚀
