# Implementation Plan: Fix CORS Configuration

**Branch**: `004-fix-cors-config` | **Date**: 2026-05-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-fix-cors-config/spec.md`

## Summary

Fix CORS configuration to allow frontend at `https://www.l9.lc` to communicate with backend API at `https://api.l9.lc`. The issue is that `cors_origins` configuration is not being applied from the database, causing the CORS middleware to reject all non-localhost origins.

**Root Cause**: The admin UI field for CORS origins exists but the configuration is not being properly saved/loaded from the database. The CORS middleware correctly checks `config.Config.Router.CORSOrigins` but this value remains empty.

**Solution**: Verify the admin UI correctly binds to the `cors_origins` field and ensure the configuration is persisted and loaded correctly.

## Technical Context

**Language/Version**: Go 1.23 (backend), Vue 3 + TypeScript (admin frontend)
**Primary Dependencies**: Fiber v2 (Go web framework), Arco Design (UI components)
**Storage**: SQLite (database for config storage)
**Testing**: Manual testing via curl and browser
**Target Platform**: Linux server (VPS) + Vercel (frontend)
**Project Type**: Web application (CMS with admin panel)
**Performance Goals**: CORS check must add <1ms latency
**Constraints**: Must work through Cloudflare CDN and Nginx reverse proxy
**Scale/Scope**: Single admin user configuring CORS, affects all API requests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

✅ No constitution violations - this is a configuration fix, not a new feature.

## Project Structure

### Documentation (this feature)

```text
specs/004-fix-cors-config/
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
├── api/web/middleware/cors.go       # CORS middleware (already exists)
├── domain/config/entity/router.go   # Router config entity (already exists)
└── domain/config/aggregate/config.go # Config aggregate

admin/
└── src/views/config/module/router/options.vue  # Admin UI (already modified)
```

**Structure decision**: Existing web application structure - no new files needed, only verification and potential fixes.

## Complexity Tracking

No complexity violations - this is a straightforward configuration fix.
