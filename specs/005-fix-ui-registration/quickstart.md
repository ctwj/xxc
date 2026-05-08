# Quickstart: Fix UI Theme Toggle and User Registration

**Feature**: Fix UI Theme Toggle and User Registration
**Date**: 2026-05-08

## Prerequisites

- Go 1.23+ installed
- Node.js 18+ and pnpm installed
- SQLite database (default)
- Access to `xxc.zip` for theme reference code

## Quick Test Scenarios

### Scenario 1: Test User Registration (After Fix)

1. Start the backend server:
   ```bash
   cd main && go run ./cmd/web/main.go
   ```

2. Test registration with curl:
   ```bash
   curl -X POST http://localhost:9008/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
   ```

3. Expected success response:
   ```json
   {
     "success": true,
     "user": {
       "id": 1,
       "username": "testuser",
       "email": "test@example.com",
       "role": "user"
     }
   }
   ```

4. Verify user table exists:
   ```bash
   sqlite3 main/moss.db "SELECT * FROM user;"
   ```

### Scenario 2: Test Theme Toggle

1. Start admin frontend:
   ```bash
   cd admin && pnpm dev
   ```

2. Open browser to `http://localhost:3000/admin/`

3. Login with admin credentials

4. Click the dark/light mode toggle (sun/moon icon in header)

5. Verify:
   - Dark mode: Body has `arco-theme="dark"` attribute, dark background
   - Light mode: Body has no `arco-theme` attribute, light background
   - All UI components render correctly in both modes

### Scenario 3: Compare with Reference Code

1. Extract `xxc.zip`:
   ```bash
   unzip xxc.zip -d xxc-reference
   ```

2. Compare theme-related files:
   ```bash
   # Compare store
   diff admin/src/store/index.js xxc-reference/admin/src/store/index.js

   # Compare base layout
   diff admin/src/layout/base.vue xxc-reference/admin/src/layout/base.vue

   # Compare dark toggle component
   diff admin/src/components/app/Dark.vue xxc-reference/admin/src/components/app/Dark.vue

   # Compare styles
   diff admin/src/style.css xxc-reference/admin/src/style.css
   ```

3. Identify and apply any missing styles or configurations

## Verification Checklist

- [ ] User table created in database
- [ ] Registration returns success with valid data
- [ ] Registration returns appropriate error for duplicate username
- [ ] Registration returns appropriate error for duplicate email
- [ ] Dark mode toggle applies dark theme correctly
- [ ] Light mode toggle applies light theme correctly
- [ ] Theme preference persists after page refresh
- [ ] All UI components render correctly in both themes
