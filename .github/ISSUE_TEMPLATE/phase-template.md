---
name: Phase Implementation
about: Template for implementing a phase of a service
title: 'Phase [N]: [Phase Name] - [Service Name]'
labels: enhancement, [service-name]
assignees: ''
---

## Phase [N]: [Phase Name]

**Service:** [service-name]-api
**Priority:** High/Medium/Low
**Estimated Time:** X hours
**Branch:** `feature/[service-name]/[issue-number]-[short-description]`

---

### Description
[Brief description of what this phase accomplishes]

---

### Tasks
- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
- [ ] Task 4

---

### Success Criteria
- [ ] Criteria 1
- [ ] Criteria 2
- [ ] Criteria 3

---

### Dependencies
**Requires:**
- Phase X completed

**Blocks:**
- Phase Y

---

### Files to Create/Modify
```
path/to/file1.go
path/to/file2.go
path/to/file3.go
```

---

### Testing Steps
```bash
# Compilation check
cd backend/[service-name]-api
go mod tidy
go build ./cmd/api

# Run tests
go test ./... -v

# Manual testing
[specific commands for this phase]
```

---

### Implementation Guide
1. Create branch: `git checkout -b feature/[service-name]-api-phase-[N]-[name]`
2. Use plan mode: `@CONTEXT_[SERVICE]_API.md Implement Phase [N]`
3. Follow success criteria
4. Test thoroughly
5. Commit: `git commit -m "feat([service]): [description]"`
6. Push and create PR to `dev`

---

### References
- Implementation Plan: See `GITFLOW.md` Phase [N]
- Context: `CONTEXT_[SERVICE]_API.md`
- Pattern Reference: `backend/users-api/[similar-file]`

---

### Notes
[Any additional notes, gotchas, or important considerations]
