# CarPooling - Development Guide

## 📚 Documentation Index

This project has structured documentation to guide development:

| Document | Purpose | When to Use |
|----------|---------|-------------|
| `CONTEXT_BOOKINGS_API.md` | Complete specification for bookings-api | Reference in plan mode for every phase |
| `GITFLOW.md` | Git workflow, branch strategy, all phases breakdown | Before starting any phase, for issue details |
| `HOW_TO_USE_PLAN_MODE.md` | Step-by-step guide for using Claude Code | When starting development on a phase |
| `.github/ISSUE_TEMPLATE/` | GitHub issue templates | When creating issues for phases |
| `README.md` | Project overview | General information |

---

## 🚀 Quick Start

### 1. First Time Setup
```bash
cd /home/user/CarPooling

# Review project structure
cat CONTEXT_BOOKINGS_API.md
cat GITFLOW.md

# Understand the workflow
cat HOW_TO_USE_PLAN_MODE.md
```

### 2. Start Development on bookings-api

#### Create GitHub Issues (Recommended)
```bash
# Go to: https://github.com/juan-cabra1/CarPooling/issues
# Click "New Issue"
# Choose template: "Bookings API - Phase 1"
# Create issues for all phases you plan to work on
```

#### Start Phase 1
```bash
# Checkout dev branch
git checkout dev
git pull origin dev

# Create feature branch
git checkout -b feature/bookings-api-phase-1-setup

# Start Claude Code
claude

# In Claude Code, type:
# @CONTEXT_BOOKINGS_API.md Implement Phase 1 (Project Setup) for bookings-api
# [Review plan, approve, let Claude implement]

# Test
cd backend/bookings-api
go mod tidy
go build ./cmd/api
go run cmd/api/main.go

# Commit
git add .
git commit -m "feat(bookings): implement phase 1 - project setup

- Initialize Go modules
- Setup .env configuration
- Create main.go skeleton
- Add all dependencies

Closes #1"

# Push and create PR
git push -u origin feature/bookings-api-phase-1-setup
# Create PR on GitHub: feature/bookings-api-phase-1-setup → dev
```

---

## 📋 Development Workflow Summary

```
┌─────────────────────────────────────────────────────────┐
│ 1. Read GITFLOW.md for current phase details           │
│ 2. Create GitHub issue from template                    │
│ 3. Create feature branch                                │
│ 4. Use plan mode with @CONTEXT_BOOKINGS_API.md         │
│ 5. Review and approve plan                             │
│ 6. Implement (Claude does this)                        │
│ 7. Test (verify success criteria)                      │
│ 8. Commit with conventional commit message             │
│ 9. Push and create PR to dev                           │
│ 10. After merge, start next phase                      │
└─────────────────────────────────────────────────────────┘
```

---

## 📁 Project Structure

```
CarPooling/
├── backend/
│   ├── users-api/          ✅ DONE - Reference implementation
│   ├── bookings-api/       🚧 TO DO - Next focus
│   ├── trips-api/          📋 TODO - After bookings
│   └── search-api/         📋 TODO - Final service
│
├── .github/
│   └── ISSUE_TEMPLATE/     GitHub issue templates
│
├── CONTEXT_BOOKINGS_API.md    📖 Full specification
├── GITFLOW.md                 📖 Workflow & phases
├── HOW_TO_USE_PLAN_MODE.md   📖 Usage guide
├── README.md                  📖 Project overview
└── README_DEVELOPMENT.md      📖 This file
```

---

## 🎯 Implementation Phases for bookings-api

### Week 1: Foundation (Phases 1-6)
- **Phase 1:** Project Setup (2-3h)
- **Phase 2:** Database Setup (3-4h)
- **Phase 3:** Repository Layer (3-4h)
- **Phase 4:** Business Logic (4-5h)
- **Phase 5:** HTTP Controllers (3-4h)
- **Phase 6:** JWT Middleware (1-2h)

### Week 2: Integration (Phases 7-12) ⚠️ CRITICAL
- **Phase 7:** RabbitMQ Consumer + Idempotency (6-8h) 🔴 MOST IMPORTANT
- **Phase 8:** RabbitMQ Publisher (3-4h)
- **Phase 9:** HTTP Clients (2-3h)
- **Phase 10:** Error Handling & Logging (2-3h)
- **Phase 11:** Routes Setup (1-2h)
- **Phase 12:** Main Assembly (2-3h)

### Week 3: Production (Phases 13-15)
- **Phase 13:** Docker & Config (2-3h)
- **Phase 14:** Testing (4-6h)
- **Phase 15:** Integration Testing (3-4h)

**Total Estimated:** 45-60 hours

---

## 🔑 Key Concepts

### Idempotency (CRITICAL)
**Problem:** RabbitMQ retries → duplicate events → double bookings

**Solution:** Check event_id before processing
```go
shouldProcess, _ := idempotencyService.CheckAndMarkEvent(event.EventID, event.EventType)
if !shouldProcess {
    logger.Info().Msg("Event already processed, skipping")
    return nil // ACK without processing
}
// Process event...
```

### Event Flow
```
1. POST /bookings → create booking (pending)
2. Publish reservation.created
3. trips-api decreases seats
4. trips-api publishes trip.updated
5. bookings-api consumes → updates booking (confirmed)
```

### Architecture Pattern
```
Controller → Service → Repository → DAO (GORM) → MySQL
                ↓
            Publisher → RabbitMQ
                          ↓
Consumer → IdempotencyService → Repository
```

---

## 🧪 Testing Checklist

After each phase:
```bash
# ✅ Compilation
go mod tidy
go build ./cmd/api

# ✅ Unit tests
go test ./... -v

# ✅ Manual testing
go run cmd/api/main.go
curl http://localhost:8003/health

# ✅ Git status clean
git status
```

---

## 📝 Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

**Types:** feat, fix, refactor, test, docs, chore

**Examples:**
```
feat(bookings): implement idempotency service
test(bookings): add race condition test for duplicate events
refactor(bookings): extract error types to domain package
docs(bookings): update API documentation
```

---

## 🐛 Common Issues & Solutions

### "Claude can't find context file"
```bash
# Use absolute path
@/home/user/CarPooling/CONTEXT_BOOKINGS_API.md

# Or paste content directly
```

### "Tests failing after implementation"
```bash
# Run specific test with verbose output
go test ./internal/service/idempotency_service_test.go -v

# Check MySQL
mysql -u root -p
USE carpooling_bookings;
SHOW TABLES;
```

### "RabbitMQ connection failed"
```bash
# Check if RabbitMQ is running
sudo systemctl status rabbitmq-server

# Or with Docker
docker ps | grep rabbit
```

### "Duplicate bookings created"
```bash
# Verify idempotency table
mysql -u root -p
USE carpooling_bookings;
SELECT * FROM processed_events;

# Check UNIQUE constraint exists
SHOW CREATE TABLE processed_events;
```

---

## 📊 Progress Tracking

Track your progress by checking off phases:

**bookings-api:**
- [ ] Phase 1: Project Setup
- [ ] Phase 2: Database Setup
- [ ] Phase 3: Repository Layer
- [ ] Phase 4: Business Logic
- [ ] Phase 5: HTTP Controllers
- [ ] Phase 6: JWT Middleware
- [ ] Phase 7: RabbitMQ Consumer (CRITICAL)
- [ ] Phase 8: RabbitMQ Publisher
- [ ] Phase 9: HTTP Clients
- [ ] Phase 10: Error Handling & Logging
- [ ] Phase 11: Routes Setup
- [ ] Phase 12: Main Assembly
- [ ] Phase 13: Docker & Config
- [ ] Phase 14: Testing
- [ ] Phase 15: Integration Testing

---

## 🎓 Learning Resources

- **Go GORM:** https://gorm.io/docs/
- **Gin Framework:** https://gin-gonic.com/docs/
- **RabbitMQ Go:** https://www.rabbitmq.com/tutorials/tutorial-two-go.html
- **Idempotency Patterns:** https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/
- **Reference Code:** `backend/users-api/` (working implementation)

---

## 👥 Team Collaboration

### Before Starting Work
1. Assign issue to yourself
2. Comment on issue: "Starting work on this"
3. Create branch from dev

### During Work
1. Commit frequently with clear messages
2. Push to your branch regularly
3. Update issue with progress/blockers

### After Completion
1. Verify all success criteria
2. Create PR with checklist from issue
3. Request review from team
4. Address review comments
5. Merge after approval

---

## 🚦 Definition of Done

A phase is DONE when:
✅ All tasks completed
✅ All success criteria met
✅ Code compiles without warnings
✅ Tests pass
✅ Manual testing successful
✅ Code committed with good message
✅ PR created and approved
✅ Merged to dev
✅ Issue closed

---

## 📞 Need Help?

1. **Check documentation:** Read CONTEXT, GITFLOW, HOW_TO_USE_PLAN_MODE
2. **Review reference:** Look at users-api implementation
3. **Search issue:** Someone may have had same problem
4. **Ask Claude:** Provide context and specific question
5. **Team discussion:** Use team communication channel

---

## 🎯 Next Steps

1. **Read all documentation:**
   - `CONTEXT_BOOKINGS_API.md`
   - `GITFLOW.md`
   - `HOW_TO_USE_PLAN_MODE.md`

2. **Create GitHub issues:**
   - Use templates in `.github/ISSUE_TEMPLATE/`
   - Create issues for phases you'll work on

3. **Start Phase 1:**
   - Follow workflow in `HOW_TO_USE_PLAN_MODE.md`
   - Reference `GITFLOW.md` for details

4. **Stay focused:**
   - One phase at a time
   - Test thoroughly
   - Commit often

---

## 📈 Success Indicators

You're doing it right if:
✅ Each phase takes 2-8 hours
✅ You understand code before committing
✅ Tests pass consistently
✅ PRs are focused and small
✅ Git history is clean
✅ Documentation stays updated

---

Happy coding! 🚀

For detailed instructions, see:
- `HOW_TO_USE_PLAN_MODE.md` - Usage guide
- `GITFLOW.md` - All phases breakdown
- `CONTEXT_BOOKINGS_API.md` - Complete specification
