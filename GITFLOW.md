# GitFlow Strategy - CarPooling Project

## Branch Strategy

### Main Branches
- `main` - Production-ready code (entrega final)
- `dev` - Integration branch (desarrollo activo)

### Feature Branches Pattern
```
feature/{service-name}/{issue-number}-{short-description}
```

**Examples:**
- `feature/trips-api/1-project-setup`
- `feature/trips-api/2-mongodb-connection`
- `feature/trips-api/7-rabbitmq-consumer`

**Why this pattern:**
- Clear service identification
- Links to GitHub issue automatically
- Short description for context
- Easy to track in git log

### Workflow
```
main (protected, no direct commits)
  └── dev (integration branch, all PRs go here)
       ├── feature/trips-api/1-project-setup
       ├── feature/trips-api/2-mongodb-connection
       ├── feature/trips-api/3-domain-models
       ├── feature/bookings-api/1-project-setup
       └── ...
```

---

## Development Workflow

### 1. Before Starting Work
```bash
# Always start from updated dev
git checkout dev
git pull origin dev

# Verify you're on latest dev
git log --oneline -5
```

### 2. Create Feature Branch
```bash
# Pattern: feature/{service}/{issue-number}-{description}
git checkout -b feature/trips-api/1-project-setup

# Verify branch created
git branch --show-current
```

### 3. Work on Feature
```bash
# Make changes following the plan
# Commit frequently (every logical change)

# Stage changes
git add .

# Commit with conventional commit message
git commit -m "feat(trips): initialize Go project structure

- Create go.mod and go.sum
- Setup basic main.go
- Add initial dependencies

Closes #1"
```

### 4. Push and Create PR
```bash
# Push feature branch
git push -u origin feature/trips-api/1-project-setup

# Create PR on GitHub:
# - Base: dev
# - Compare: feature/trips-api/1-project-setup
# - Title: "feat(trips): Project setup - Issue #1"
# - Link issue in description
```

### 5. After PR Merged
```bash
# Update dev branch
git checkout dev
git pull origin dev

# Delete local feature branch (cleanup)
git branch -d feature/trips-api/1-project-setup

# Start next feature
git checkout -b feature/trips-api/2-mongodb-connection
```

---

## Commit Message Convention

### Format
```
type(scope): short description

[optional body explaining the change]

[optional footer with issue references]
```

### Types
- `feat` - New feature
- `fix` - Bug fix
- `refactor` - Code refactoring (no feature change)
- `test` - Adding or updating tests
- `docs` - Documentation only
- `chore` - Build, dependencies, config changes
- `style` - Code style changes (formatting)

### Scopes (by service)
- `trips` - trips-api changes
- `bookings` - bookings-api changes
- `users` - users-api changes
- `search` - search-api changes
- `project` - Root project changes

### Examples
```bash
# Simple feature
feat(trips): add trip repository with MongoDB

# Feature with details
feat(trips): implement RabbitMQ consumer for reservation events

- Add consumer for reservation.created
- Add consumer for reservation.cancelled
- Implement idempotency check

Closes #7

# Bug fix
fix(trips): correct available seats calculation

# Refactoring
refactor(trips): extract domain models to separate package

# Tests
test(trips): add idempotency service unit tests

# Documentation
docs: update README with trips-api setup instructions

# Dependencies
chore(trips): add zerolog for structured logging
```

### Rules
- ✅ Use lowercase for type and scope
- ✅ Keep first line under 72 characters
- ✅ Use present tense ("add" not "added")
- ✅ Reference issue number in footer
- ✅ Explain "why" in body, not "what" (code shows what)

---

## Issues Structure

### Creating Issues

**One issue per feature/phase:**
- Each issue = One logical piece of work
- Small enough to complete in 1-2 days
- Large enough to be meaningful

**Issue Naming Convention:**
```
[Service] Feature Name - Brief Description
```

**Examples:**
- `[trips-api] Project Setup`
- `[trips-api] MongoDB Connection & Models`
- `[trips-api] RabbitMQ Consumer with Idempotency`
- `[bookings-api] Repository Layer Implementation`

### Issue Labels

Use GitHub labels for organization:
- `trips-api`, `bookings-api`, `users-api`, `search-api` - Service identification
- `enhancement` - New feature
- `bug` - Bug fix
- `documentation` - Documentation work
- `critical` - Must be done / high priority
- `phase-N` - Phase number (optional)

### Issue Template Structure

Will be provided in `.github/ISSUE_TEMPLATE/` with details for:
- Description of the work
- Tasks checklist
- Success criteria
- Files to modify
- Testing steps
- Branch name to use

---

## PR Review Checklist

Before approving any PR, verify:
- [ ] Code compiles without warnings
- [ ] All tests pass (if applicable)
- [ ] No hardcoded values (use config/env)
- [ ] Error handling is proper
- [ ] Logging is appropriate
- [ ] No sensitive data in commits (.env files gitignored)
- [ ] README or docs updated if needed
- [ ] Success criteria from issue met
- [ ] Branch is up to date with dev

---

## Service Implementation Order

1. **trips-api** (MongoDB) - Main API, implement first
2. **bookings-api** (MySQL) - Depends on trips-api
3. **search-api** (MongoDB + Solr) - Depends on trips-api events
4. **users-api** - ✅ Already complete

---

## Next Steps

1. **Read service-specific CONTEXT files** for detailed specifications
2. **Create GitHub Issues** for the service you're working on
3. **Follow the workflow** above for each issue
4. **Use plan mode** with Claude Code for implementation
5. **Test thoroughly** before creating PR
6. **Get review** and merge to dev
