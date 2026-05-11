# Tasks: Fix Content Display Issues

**Input**: Design documents from `/specs/009-fix-content-display/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Foundational (Shared Utility)

**Purpose**: Create the HTML tag stripping utility used by both US1 tasks

- [x] T001 Add `stripHTMLTags` helper function using `regexp` to remove HTML tags from strings in `main/plugins/TelegramChannelSync.go`

**Checkpoint**: Helper function ready for use by US1 tasks

---

## Phase 2: User Story 1 - Fix HTML Tags in Article Descriptions (Priority: P1) MVP

**Goal**: Article card descriptions on the homepage display clean plain text without any visible HTML tags

**Independent Test**: Open homepage, verify all 23 article card descriptions contain no `<p>`, `</p>`, `<br>`, `<strong>` tags

### Implementation for User Story 1

- [x] T002 [US1] Modify `truncateDescription()` in `main/plugins/TelegramChannelSync.go` to call `stripHTMLTags()` before truncating content to 200 runes — fixes new articles at storage time
- [x] T003 [P] [US1] Strip HTML tags from `description` field in `APIArticleList()` and `APIArticleDetail()` responses in `main/api/web/controller/api.go` — fixes existing articles in database

**Checkpoint**: Homepage cards show clean text descriptions; no HTML tags visible anywhere

---

## Phase 3: User Story 2 - Fix Article Images Not Displaying (Priority: P1)

**Goal**: All article thumbnail images and inline images load correctly on homepage and detail pages

**Independent Test**: Open homepage, verify thumbnail images load on all cards with media; open an article detail page, verify inline images display

### Implementation for User Story 2

- [x] T004 [P] [US2] Add `GET /telegram/media/:mediaId` route to `RegisterAPIRoutes()` in `main/api/web/router/api.go`, pointing to `controller.TelegramGetMedia` — makes media accessible at `/api/telegram/media/{id}` matching stored URLs

**Checkpoint**: All article images load successfully (HTTP 200) at `/api/telegram/media/{id}`

---

## Phase 4: User Story 3 - Fix Article Videos Not Displaying (Priority: P2)

**Goal**: Video content in articles plays correctly

**Independent Test**: Open the article with video content, verify video player loads and video plays

### Verification for User Story 3

- [x] T005 [US3] Verify video loads correctly via `/api/telegram/media/{id}` — same route fix from T004 covers both images and videos; confirm by testing the article with `videoUrl` field

**Checkpoint**: Article video plays without errors

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [x] T006 Verify all acceptance scenarios from spec.md pass: (1) no HTML tags on cards, (2) images load, (3) videos play, (4) media accessible without auth, (5) admin panel unaffected
- [x] T007 Build and smoke test: run `task build` and verify the production binary serves correct API responses

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: No dependencies — start immediately
- **US1 (Phase 2)**: T002 depends on T001 (stripHTMLTags); T003 depends on T001 (stripHTMLTags)
- **US2 (Phase 3)**: Independent of US1 — can start immediately (different files)
- **US3 (Phase 4)**: Depends on T004 (same route)
- **Polish (Phase 5)**: Depends on all user stories being complete

### Task Dependencies

```
T001 ──→ T002 (stripHTMLTags used in truncateDescription)
T001 ──→ T003 (stripHTMLTags used in API controller)
T004 ──→ T005 (same route, T005 just verifies)
T002 + T003 + T004 + T005 ──→ T006 + T007 (validation)
```

### Parallel Opportunities

- **T003 and T004** can run in parallel (different files: `controller/api.go` vs `router/api.go`)
- **T002 and T004** can run in parallel (different files: `TelegramChannelSync.go` vs `router/api.go`)

### Parallel Example: After T001

```bash
# Launch these in parallel (all touch different files):
Task T002: "Modify truncateDescription() in TelegramChannelSync.go"
Task T003: "Strip HTML in controller/api.go"
Task T004: "Add media route in router/api.go"
```

---

## Implementation Strategy

### MVP First (US1 + US2)

1. Complete Phase 1: Add `stripHTMLTags` helper (T001)
2. Complete Phase 2 + Phase 3 in parallel: Fix HTML tags (T002+T003) AND Fix images (T004)
3. **STOP and VALIDATE**: Test both fixes on dev environment
4. Deploy if ready — this resolves all critical user-facing issues

### Full Delivery

1. MVP above
2. Phase 4: Verify videos (T005)
3. Phase 5: Full validation and build (T006+T007)

---

## Notes

- Total tasks: 7 (4 implementation + 1 verification + 2 validation)
- Files modified: 3 Go files only
- No frontend changes needed
- No database migration needed
- Existing stored media URLs (`/api/telegram/media/{id}`) work without data changes
- Admin panel media route (`/admin/api/telegram/media/{id}`) remains unaffected
