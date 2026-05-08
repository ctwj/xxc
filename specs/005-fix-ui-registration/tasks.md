# Tasks: Fix UI Theme Toggle and User Registration

**Input**: Design documents from `/specs/005-fix-ui-registration/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Manual testing only - no automated tests requested.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `main/` at repository root
- **Admin Frontend**: `admin/` at repository root

---

## Phase 1: Setup (Verification)

**Purpose**: Verify existing code structure and identify what's already in place

- [x] T001 Verify User entity exists in `main/domain/core/entity/user.go`
- [x] T002 Verify UserRepository exists in `main/domain/core/repository/user.go`
- [x] T003 [P] Verify theme store state in `admin/src/store/index.js`
- [x] T004 [P] Verify theme toggle component in `admin/src/components/app/Dark.vue`
- [x] T005 [P] Verify theme attribute management in `admin/src/layout/base.vue`

---

## Phase 2: Foundational (Database Migration Fix)

**Purpose**: Fix the missing user table migration - this MUST be complete before registration can work

**⚠️ CRITICAL**: This is the root cause of the registration failure

- [x] T006 Add `MigrateTable()` method to `UserRepository` in `main/domain/core/repository/user.go`
- [x] T007 Add `User.MigrateTable()` call to `repository.MigrateTable()` in `main/domain/core/repository/repository.go`

**Checkpoint**: User table will be created on next application start

---

## Phase 3: User Story 1 - Dark/Light Mode Toggle Works Correctly (Priority: P1) 🎯 MVP

**Goal**: Fix theme toggle so UI displays correctly in both dark and light modes

**Independent Test**: Toggle the dark/light mode switch in admin panel and verify all UI elements display correctly in both modes

### Implementation for User Story 1

- [x] T008 [US1] Extract `xxc.zip` to temporary directory for comparison
- [x] T009 [US1] Compare `admin/src/store/index.js` with reference code
- [x] T010 [P] [US1] Compare `admin/src/layout/base.vue` with reference code
- [x] T011 [P] [US1] Compare `admin/src/components/app/Dark.vue` with reference code
- [x] T012 [P] [US1] Compare `admin/src/style.css` with reference code
- [x] T013 [US1] Apply any missing styles or configurations identified from comparison
- [ ] T014 [US1] Test theme toggle in browser - verify dark mode displays correctly
- [ ] T015 [US1] Test theme toggle in browser - verify light mode displays correctly
- [ ] T016 [US1] Test theme persistence - refresh page and verify theme state persists

**Checkpoint**: At this point, User Story 1 should be fully functional - theme toggle works correctly

---

## Phase 4: User Story 2 - User Registration Works After Deployment (Priority: P1) 🎯 MVP

**Goal**: Fix user registration so new users can create accounts

**Independent Test**: Visit frontend registration page, enter valid credentials, and verify successful account creation

### Implementation for User Story 2

- [x] T017 [US2] Rebuild backend with migration fix: `cd main && go build ./cmd/web/main.go`
- [x] T018 [US2] Stop the running Moss service on server
- [x] T019 [US2] Deploy new binary to server
- [x] T020 [US2] Start the Moss service - this will run auto-migration and create user table
- [x] T021 [US2] Test registration with curl:
  ```bash
  curl -X POST http://localhost:9008/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
  ```
- [x] T022 [US2] Verify success response with user data
- [x] T023 [US2] Test duplicate username error handling
- [x] T024 [US2] Test duplicate email error handling
- [x] T025 [US2] Verify user table exists in database: `sqlite3 moss.db "SELECT * FROM user;"`

**Checkpoint**: At this point, User Story 2 should be fully functional - users can register successfully

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and cleanup

- [x] T026 [P] Update quickstart.md with verified test steps
- [x] T027 [P] Document theme toggle fix in deployment notes
- [x] T028 Commit all changes with descriptive message

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - verification only
- **Foundational (Phase 2)**: No dependencies - can start immediately after verification
- **User Story 1 (Phase 3)**: No dependencies on Phase 2 - can run in parallel
- **User Story 2 (Phase 4)**: Depends on Phase 2 (migration fix must be implemented first)
- **Polish (Phase 5)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (Theme Toggle)**: Independent - can be worked on in parallel with User Story 2
- **User Story 2 (Registration)**: Depends on Foundational (Phase 2) being complete

### Within Each User Story

- Verification before implementation
- Code changes before testing
- Test before moving to next story

### Parallel Opportunities

- T003, T004, T005 can run in parallel (different files)
- T010, T011, T012 can run in parallel (different files)
- T026, T027 can run in parallel (different files)
- User Story 1 and User Story 2 can be worked on in parallel after Phase 2

---

## Parallel Example: Setup Phase

```bash
# Launch verification tasks in parallel:
Task: "Verify User entity in main/domain/core/entity/user.go"
Task: "Verify UserRepository in main/domain/core/repository/user.go"
Task: "Verify theme store in admin/src/store/index.js"
Task: "Verify theme toggle in admin/src/components/app/Dark.vue"
Task: "Verify theme management in admin/src/layout/base.vue"
```

---

## Implementation Strategy

### MVP First (Both User Stories are P1)

1. Complete Phase 1: Setup (verification)
2. Complete Phase 2: Foundational (migration fix) - CRITICAL for registration
3. Complete Phase 3: User Story 1 (theme toggle) - Can run in parallel with Phase 2
4. Complete Phase 4: User Story 2 (registration) - Requires Phase 2 first
5. **STOP and VALIDATE**: Test both user stories independently
6. Deploy if ready

### Incremental Delivery

1. Complete Setup → Verification confirms code structure
2. Add Foundational → Migration fix enables registration
3. Add User Story 1 → Test independently → Theme toggle works
4. Add User Story 2 → Test independently → Registration works
5. Each story adds value without breaking previous stories

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- This is primarily a bug fix task, not new feature development
- Focus is on fixing existing code, not creating new files
- The root cause of registration failure is missing database migration
- The theme issue requires comparison with reference code in `xxc.zip`
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
