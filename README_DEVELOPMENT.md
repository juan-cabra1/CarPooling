# CarPooling - Development Guide

## 📚 Documentation Index

| Document | Purpose | When to Use |
|----------|---------|-------------|
| `README_DEVELOPMENT.md` | **This file** - Quick start and navigation | Start here |
| `GITFLOW.md` | Git workflow, branch strategy, commit conventions | Before creating any branch |
| `CONTEXT_TRIPS_API.md` | Complete specification for trips-api | When working on trips-api |
| `CONTEXT_BOOKINGS_API.md` | Complete specification for bookings-api | When working on bookings-api |
| `HOW_TO_USE_PLAN_MODE.md` | Step-by-step guide for using Claude Code | When starting a new feature |
| `.github/ISSUE_TEMPLATE/` | GitHub issue templates | When creating issues |

---

## 🚀 Quick Start

### First Time Setup

```bash
# 1. Clone repository (if not done)
cd /home/user/CarPooling

# 2. Read documentation
cat README_DEVELOPMENT.md  # This file
cat GITFLOW.md             # Git workflow

# 3. Choose service to work on
# - trips-api (start here - main API)
# - bookings-api (after trips-api)
# - search-api (after trips-api)
# - users-api (already done ✅)
```

### Start Working on a Feature

```bash
# 1. Update dev branch
git checkout dev
git pull origin dev

# 2. Read service context
cat CONTEXT_TRIPS_API.md  # Or CONTEXT_BOOKINGS_API.md

# 3. Create GitHub issue for feature
# Go to: https://github.com/juan-cabra1/CarPooling/issues
# Use templates in .github/ISSUE_TEMPLATE/

# 4. Create feature branch
git checkout -b feature/trips-api/1-project-setup

# 5. Use Claude Code plan mode
claude
# Then: @CONTEXT_TRIPS_API.md Implement feature X...

# 6. Test, commit, push
# See GITFLOW.md for details
```

---

## 📋 Project Structure

```
CarPooling/
├── backend/
│   ├── users-api/          ✅ DONE (MySQL, Gin, GORM)
│   ├── trips-api/          🚧 TO DO - Start here (MongoDB)
│   ├── bookings-api/       📋 TODO - After trips-api (MySQL)
│   └── search-api/         📋 TODO - After trips-api (MongoDB + Solr)
│
├── .github/
│   └── ISSUE_TEMPLATE/     Issue templates for features
│
├── CONTEXT_TRIPS_API.md       Trips API specification
├── CONTEXT_BOOKINGS_API.md    Bookings API specification
├── GITFLOW.md                 Git workflow & conventions
├── HOW_TO_USE_PLAN_MODE.md    Claude Code usage guide
├── README.md                  Project overview
└── README_DEVELOPMENT.md      This file
```

---

## 🎯 Development Workflow (Summary)

1. **Pick a feature** → Read issue or create one
2. **Create branch** → `feature/{service}/{issue-number}-{description}`
3. **Read context** → `CONTEXT_{SERVICE}_API.md`
4. **Use plan mode** → Claude Code with @context
5. **Implement** → Follow the plan
6. **Test** → Compile, run tests, manual testing
7. **Commit** → Conventional commit message
8. **Push & PR** → Create PR to `dev`
9. **Review** → Wait for approval
10. **Merge & repeat** → Start next feature

**Details:** See `GITFLOW.md` for complete workflow

---

## 🏗️ Service Implementation Order

### 1. trips-api (Start Here) 🎯
- **Database:** MongoDB
- **Port:** 8002
- **Why first:** Main API, other services depend on it
- **Context:** `CONTEXT_TRIPS_API.md`

### 2. bookings-api (After trips-api)
- **Database:** MySQL
- **Port:** 8003
- **Depends on:** trips-api (validates trips, consumes events)
- **Context:** `CONTEXT_BOOKINGS_API.md`

### 3. search-api (After trips-api)
- **Database:** MongoDB + Solr
- **Port:** 8004
- **Depends on:** trips-api (consumes events for indexing)
- **Context:** TBD

### 4. users-api ✅
- **Already complete** - Use as reference

---

## 🔑 Key Concepts

### Idempotency (CRITICAL)
Services consume RabbitMQ events. If RabbitMQ retries a message, we must not process it twice.

**Solution:** Check `event_id` before processing
- trips-api: Uses MongoDB unique index on `event_id`
- bookings-api: Uses MySQL unique constraint on `event_id`

### Event-Driven Architecture
```
trips-api publishes → RabbitMQ → bookings-api/search-api consume
bookings-api publishes → RabbitMQ → trips-api consumes
```

### Service Communication
- **Synchronous:** HTTP REST (e.g., bookings calls trips-api to validate)
- **Asynchronous:** RabbitMQ events (e.g., trip.updated, reservation.created)

---

## 🧪 Testing Checklist

After each feature:
```bash
# ✅ Compilation
go mod tidy
go build ./cmd/api

# ✅ Unit tests (if applicable)
go test ./... -v

# ✅ Manual testing
go run cmd/api/main.go
curl http://localhost:800X/health

# ✅ Git status clean
git status
git diff
```

---

## 📝 Commit Message Format

```
type(scope): short description

[optional body]

Closes #issue-number
```

**Examples:**
```
feat(trips): add trip repository with MongoDB driver
fix(bookings): correct seat availability validation
docs: update CONTEXT_TRIPS_API.md with event schema
```

**See:** `GITFLOW.md` for complete convention

---

## 🐛 Common Issues

### "Can't connect to MongoDB"
```bash
# Check if MongoDB is running
docker ps | grep mongo

# Or start with docker-compose
cd backend/search-api  # Has docker-compose with mongo
docker-compose up -d mongo
```

### "Can't connect to RabbitMQ"
```bash
# Check if RabbitMQ is running
docker ps | grep rabbit

# Start RabbitMQ
docker-compose up -d rabbit
```

### "Claude Code can't find context file"
```bash
# Use absolute path in plan mode
@/home/user/CarPooling/CONTEXT_TRIPS_API.md

# Or copy content to prompt
```

---

## 📊 Progress Tracking

Use GitHub issues to track progress:
- [ ] trips-api features (see CONTEXT_TRIPS_API.md)
- [ ] bookings-api features (see CONTEXT_BOOKINGS_API.md)
- [ ] search-api features
- [x] users-api (complete)

---

## 🎓 Learning Resources

- **Go:** https://go.dev/tour/
- **MongoDB Go Driver:** https://www.mongodb.com/docs/drivers/go/current/
- **Gin Framework:** https://gin-gonic.com/docs/
- **RabbitMQ Go:** https://www.rabbitmq.com/tutorials/tutorial-two-go.html
- **GORM (for bookings-api):** https://gorm.io/docs/
- **Reference Code:** `backend/users-api/` (working implementation)

---

## 📞 Need Help?

1. **Check documentation:**
   - Read relevant CONTEXT file
   - Review GITFLOW.md for workflow
   - Check HOW_TO_USE_PLAN_MODE.md

2. **Look at users-api:**
   - Already complete and working
   - Good patterns to follow

3. **Ask Claude Code:**
   - Provide context with @filename
   - Be specific about the issue

---

## ✅ Next Steps

1. ✅ Read this file (you're here!)
2. 📖 Read `GITFLOW.md` - Understand workflow
3. 📖 Read `CONTEXT_TRIPS_API.md` - trips-api specification
4. 📖 Read `HOW_TO_USE_PLAN_MODE.md` - Claude Code usage
5. 🚀 Create first issue for trips-api
6. 🚀 Start implementing!

---

**Remember:** Small commits, frequent PRs, one feature at a time.

Happy coding! 🚀
