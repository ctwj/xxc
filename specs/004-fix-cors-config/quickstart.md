# Quickstart: Fix CORS Configuration

**Date**: 2026-05-08
**Feature**: 004-fix-cors-config

## Immediate Fix (Database Update)

If you need to fix CORS immediately without rebuilding the admin frontend:

### Step 1: Access Database

```bash
sqlite3 /opt/moss/runtime/moss.db
```

### Step 2: Check Current Config

```sql
SELECT id, data FROM config WHERE id = 'router';
```

### Step 3: Update CORS Origins

If the `data` column is JSON, you need to add `cors_origins` field:

```sql
-- Example: if current data is '{"admin_path":"/admin",...}'
-- You need to manually edit the JSON to add cors_origins
-- This is risky, so backup first!

-- Backup
CREATE TABLE config_backup AS SELECT * FROM config WHERE id = 'router';

-- Then use a JSON update approach or export/import
```

**Recommended**: Use the admin panel after rebuilding the frontend.

## Full Fix (Admin Panel)

### Step 1: Rebuild Admin Frontend

```bash
cd admin
pnpm install
pnpm run build
```

The build output goes to `main/resources/admin/`.

### Step 2: Rebuild Backend

```bash
cd main
go build -o moss ./cmd/web/main.go
```

### Step 3: Deploy

Upload the new binary to the server and restart the service.

### Step 4: Configure CORS via Admin Panel

1. Login to `https://api.l9.lc/admin/`
2. Navigate to **Settings** → **Router** → **Options**
3. Enter CORS Origins: `https://www.l9.lc`
4. Click **Save**
5. Restart the service (or the config should reload automatically)

### Step 5: Verify

```bash
curl -I -X OPTIONS "http://127.0.0.1:9008/api/auth/login" \
  -H "Origin: https://www.l9.lc" \
  -H "Access-Control-Request-Method: POST"
```

Expected response headers:
```
Access-Control-Allow-Origin: https://www.l9.lc
Access-Control-Allow-Credentials: true
```

## Verification Checklist

- [ ] CORS Origins field visible in admin panel
- [ ] Can save CORS Origins value
- [ ] Saved value persists after page refresh
- [ ] OPTIONS request returns CORS headers for allowed origin
- [ ] OPTIONS request does NOT return CORS headers for disallowed origin
- [ ] Frontend login works without CORS errors