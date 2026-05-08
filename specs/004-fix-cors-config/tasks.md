# Tasks: Fix CORS Configuration

**Input**: Design documents from `/specs/004-fix-cors-config/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Manual testing only - no automated tests requested.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `main/` at repository root
- **Admin Frontend**: `admin/` at repository root
- **Frontend (Next.js)**: `frontend/` at repository root

---

## Phase 1: Setup (Verification)

**Purpose**: Verify existing code structure and identify what's already in place

- [x] T001 Verify CORS middleware exists and is correctly implemented in `main/api/web/middleware/cors.go`
- [x] T002 Verify Router config entity has `CORSOrigins` field in `main/domain/config/entity/router.go`
- [x] T003 [P] Verify admin UI CORS Origins field exists in `admin/src/views/config/module/router/options.vue`
- [ ] T004 [P] Check current database config for `router` entry to see if `cors_origins` field exists

---

## Phase 2: Foundational (No Blocking Prerequisites)

**Purpose**: This is a configuration fix - no new foundational infrastructure needed

**⚠️ NOTE**: The existing code structure is correct. The issue is that:
1. Admin UI field was added but frontend not rebuilt
2. Configuration not saved to database

No foundational tasks required - proceed directly to user story implementation.

---

## Phase 3: User Story 1 - Configure CORS Origins via Admin Panel (Priority: P1) 🎯 MVP

**Goal**: Enable administrators to configure CORS allowed origins through the admin panel

**Independent Test**: Login to admin panel, navigate to router settings, add CORS origin, save, and verify setting persists

### Implementation for User Story 1

- [ ] T005 [US1] Rebuild admin frontend with CORS Origins field in `admin/` (pnpm run build)
- [ ] T006 [US1] Verify build output in `main/resources/admin/`
- [ ] T007 [US1] Rebuild backend binary with embedded admin assets
- [ ] T008 [US1] Deploy new binary to server
- [ ] T009 [US1] Login to admin panel at `https://api.l9.lc/admin/`
- [ ] T010 [US1] Navigate to Settings → Router → Options, enter `https://www.l9.lc` in CORS Origins field
- [ ] T011 [US1] Save and verify setting persists after page refresh

**Checkpoint**: At this point, User Story 1 should be fully functional - CORS Origins can be configured via admin panel

---

## Phase 4: User Story 2 - Verify CORS Headers in API Response (Priority: P2)

**Goal**: Verify that CORS headers are correctly returned for allowed origins

**Independent Test**: Make curl request with Origin header and inspect response headers

### Implementation for User Story 2

- [ ] T012 [US2] Test OPTIONS request with allowed origin via curl on server:
  ```bash
  curl -I -X OPTIONS "http://127.0.0.1:9008/api/auth/login" \
    -H "Origin: https://www.l9.lc" \
    -H "Access-Control-Request-Method: POST"
  ```
- [ ] T013 [US2] Verify response includes `Access-Control-Allow-Origin: https://www.l9.lc`
- [ ] T014 [US2] Test OPTIONS request with disallowed origin (e.g., `https://evil.com`)
- [ ] T015 [US2] Verify response does NOT include CORS headers for disallowed origin

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - CORS headers correctly returned

---

## Phase 5: User Story 3 - Frontend Login Works After CORS Fix (Priority: P3)

**Goal**: End users can log into the application from the frontend

**Independent Test**: Visit https://www.l9.lc, enter credentials, and successfully log in

### Implementation for User Story 3

- [ ] T016 [US3] Visit `https://www.l9.lc` in browser
- [ ] T017 [US3] Open browser developer tools → Network tab
- [ ] T018 [US3] Enter login credentials and submit
- [ ] T019 [US3] Verify login request succeeds without CORS errors
- [ ] T020 [US3] Verify session persists after page refresh

**Checkpoint**: All user stories complete - frontend can communicate with backend

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and cleanup

- [ ] T021 [P] Update quickstart.md with verified steps
- [ ] T022 [P] Document CORS configuration in deployment guide
- [ ] T023 Commit all changes with descriptive message

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - verification only
- **Foundational (Phase 2)**: N/A - no foundational tasks
- **User Story 1 (Phase 3)**: Depends on verification confirming code is correct
- **User Story 2 (Phase 4)**: Depends on User Story 1 (needs CORS configured)
- **User Story 3 (Phase 5)**: Depends on User Story 2 (needs CORS headers working)
- **Polish (Phase 6)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Setup verification
- **User Story 2 (P2)**: Depends on User Story 1 - needs CORS configured first
- **User Story 3 (P3)**: Depends on User Story 2 - needs CORS headers working first

### Within Each User Story

- Verification before implementation
- Build before deploy
- Deploy before test
- Test before moving to next story

### Parallel Opportunities

- T001, T002 can run in parallel (different files)
- T003, T004 can run in parallel (different files)
- T021, T022 can run in parallel (different files)

---

## Parallel Example: Setup Phase

```bash
# Launch verification tasks in parallel:
Task: "Verify CORS middleware in main/api/web/middleware/cors.go"
Task: "Verify Router config entity in main/domain/config/entity/router.go"
Task: "Verify admin UI field in admin/src/views/config/module/router/options.vue"
Task: "Check database config for router entry"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (verification)
2. Complete Phase 3: User Story 1 (configure CORS via admin panel)
3. **STOP and VALIDATE**: Test that CORS can be configured
4. Deploy if ready

### Incremental Delivery

1. Complete Setup → Verification confirms code is correct
2. Add User Story 1 → Test independently → Deploy (MVP!)
3. Add User Story 2 → Test independently → Verify CORS headers
4. Add User Story 3 → Test independently → End-to-end login works
5. Each story adds value without breaking previous stories

---

## Immediate Fix Alternative

If you need CORS working immediately without rebuilding the admin frontend:

### Direct Database Update

1. Access the database on the server:
   ```bash
   sqlite3 /opt/moss/runtime/moss.db
   ```

2. Check current router config:
   ```sql
   SELECT id, data FROM config WHERE id = 'router';
   ```

3. Update the JSON to add `cors_origins`:
   ```sql
   -- This requires manually editing the JSON blob
   -- Example using json_set if SQLite has JSON extension:
   UPDATE config SET data = json_set(data, '$.cors_origins', 'https://www.l9.lc') WHERE id = 'router';
   ```

4. Restart the Moss service

5. Verify with curl test

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- This is primarily a deployment/verification task, not new code development
- The code changes were already made (admin UI field added)
- Focus is on rebuilding, deploying, and verifying
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
