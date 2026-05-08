# Implementation Plan: Fix UI Theme Toggle and User Registration

**Branch**: `005-fix-ui-registration` | **Date**: 2026-05-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-fix-ui-registration/spec.md`

## Summary

Fix two issues: (1) Dark/light mode toggle not displaying correctly when switching themes in admin panel, and (2) User registration failing with "failed to create user" error after deployment.

**Root Cause Analysis**:
- **Theme Toggle**: The `store.dark` state is managed correctly in `admin/src/store/index.js` and the `arco-theme` attribute is set/removed in `admin/src/layout/base.vue`. The issue may be related to CSS not applying correctly or component-level styling issues. Reference code in `xxc.zip` for comparison.
- **Registration**: The error "failed to create user" comes from `main/domain/core/service/user.go:56` when `repository.User.Create(user)` fails. This could be due to database migration not running (user table doesn't exist) or database connection issues.

## Technical Context

**Language/Version**: Go 1.23 (backend), Vue 3 + TypeScript (admin frontend)
**Primary Dependencies**: Fiber v2 (Go web framework), Arco Design Vue (UI components), Pinia (state management)
**Storage**: SQLite (database for user storage)
**Testing**: Manual testing via browser (theme), curl/API testing (registration)
**Target Platform**: Linux server (VPS) + Vercel (frontend)
**Project Type**: Web application (CMS with admin panel)
**Performance Goals**: Theme toggle < 100ms, Registration < 2s
**Constraints**: Must work through Cloudflare CDN and Nginx reverse proxy
**Scale/Scope**: Single admin user configuring theme, affects all new user registrations

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

✅ No constitution violations - this is a bug fix, not a new feature.

## Project Structure

### Documentation (this feature)

```text
specs/005-fix-ui-registration/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
main/
├── domain/core/
│   ├── entity/user.go        # User entity (already exists)
│   ├── service/user.go       # User service with Register method (already exists)
│   └── repository/user.go    # User repository (already exists)
├── api/web/
│   ├── controller/auth.go    # Auth controller with registration endpoint
│   └── router/api.go         # API routes
└── infrastructure/persistent/db/  # Database connection

admin/
├── src/
│   ├── store/index.js        # Pinia store with dark mode state
│   ├── layout/base.vue       # Theme attribute management
│   ├── components/app/Dark.vue  # Dark mode toggle component
│   └── style.css             # Global styles including dark theme
└── package.json
```

**Structure Decision**: Existing web application structure - no new files needed, only bug fixes.

## Complexity Tracking

No complexity violations - this is a straightforward bug fix.
