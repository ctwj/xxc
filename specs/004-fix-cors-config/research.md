# Research: Fix CORS Configuration

**Date**: 2026-05-08
**Feature**: 004-fix-cors-config

## Problem Statement

The CORS middleware is not returning CORS headers for allowed origins because `config.Config.Router.CORSOrigins` is empty. The admin UI has a CORS Origins input field, but the configuration is not being persisted or loaded correctly.

## Research Findings

### 1. Configuration Save Flow

**Flow**: Admin UI → API Controller → Config Service → Repository → Database

```
admin/src/views/config/index.vue
  → configPost(id, data.value)
  → POST /admin/api/config/:id
  → controller.ConfigUpdate
  → config.Config.Save(id, ctx.Body())
  → service.Save(item, data)
  → repository.Save(item.ConfigID(), data)
  → Database (config table)
```

### 2. Configuration Load Flow

**Flow**: Database → Repository → Config Service → Config Aggregate → Admin UI

```
Database (config table)
  → repository.Get(id)
  → service.Pull(item)
  → service.Merge(data, item)
  → config.Config.Router (in memory)
```

### 3. Key Files Analysis

| File | Purpose | Status |
|------|---------|--------|
| `main/domain/config/entity/router.go` | Defines `CORSOrigins` field with JSON tag `cors_origins` | ✅ Correct |
| `main/api/web/middleware/cors.go` | Reads `config.Config.Router.CORSOrigins` | ✅ Correct |
| `admin/src/views/config/module/router/options.vue` | UI for CORS Origins input | ✅ Added |
| `main/domain/config/service/service.go` | Handles save/load with JSON marshal/unmarshal | ✅ Correct |

### 4. Root Cause Analysis

The configuration flow is correct. The issue is likely one of:

1. **Admin UI not rebuilt**: The new CORS Origins field was added to `options.vue` but the admin frontend hasn't been rebuilt and deployed.

2. **Database migration**: The `cors_origins` field exists in the entity but may not be in the database JSON blob yet.

3. **Config not reloaded**: After saving, the config may not be reloaded into memory.

### 5. Verification Steps

To verify the fix works:

1. **Check database**: Query the `config` table for `router` entry
   ```sql
   SELECT id, data FROM config WHERE id = 'router';
   ```

2. **Check if `cors_origins` is in the JSON**: The `data` column should contain `cors_origins` field after saving.

3. **Check middleware logs**: The CORS middleware logs `[CORS] origin=xxx, allowed=xxx` - check if it shows the configured origin.

### 6. Solution

**Immediate Fix**: Manually update the database to add `cors_origins` to the router config.

**Long-term Fix**: Rebuild and redeploy the admin frontend with the new CORS Origins field.

## Decisions

| Decision | Rationale | Alternatives Considered |
|----------|-----------|------------------------|
| Use existing config system | Already has JSON serialization, database persistence, and UI infrastructure | Creating a separate CORS config table - rejected as over-engineering |
| Add field to existing entity | Minimal change, follows existing patterns | Creating new entity - rejected as unnecessary complexity |

## Implementation Notes

1. The `cors_origins` field already exists in `router.go` entity
2. The admin UI field has been added to `options.vue`
3. The CORS middleware correctly reads the config
4. Need to verify the admin frontend is rebuilt and deployed
