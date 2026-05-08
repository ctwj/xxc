# Feature Specification: Fix CORS Configuration

**Feature Branch**: `004-fix-cors-config`
**Created**: 2026-05-08
**Status**: Draft
**Input**: User description: "当前系统已经成功部署完成了， 前后端部署正常， 前端显示正常， 但是接口遇到了 CORS 的错误， 请重新梳理，后端接口逻辑，已经后端配置， 提出解决方案"

## Problem Analysis

### Current Situation
- Frontend deployed at: `https://www.l9.lc` (Vercel)
- Backend API deployed at: `https://api.l9.lc` (VPS with Nginx reverse proxy)
- Backend runs on: `localhost:9008` (Go/Fiber)
- Requests go through: Cloudflare CDN → Nginx → Backend

### Root Cause
Testing on the server shows:
```bash
curl -I -X OPTIONS "http://127.0.0.1:9008/api/auth/login" \
  -H "Origin: https://www.l9.lc" \
  -H "Access-Control-Request-Method: POST"
```
Returns `HTTP/1.1 204 No Content` but **missing CORS headers**:
- No `Access-Control-Allow-Origin`
- No `Access-Control-Allow-Credentials`

This indicates the CORS middleware is not recognizing the origin as allowed.

### Current Implementation Analysis

1. **Backend CORS Middleware** (`main/api/web/middleware/cors.go`):
   - Reads `config.Config.Router.CORSOrigins` for allowed origins
   - If empty, only allows `localhost` origins
   - Correctly sets headers when origin is allowed

2. **Config Entity** (`main/domain/config/entity/router.go`):
   - Has `CORSOrigins string` field with JSON tag `cors_origins`
   - Config is stored in database, not `conf.toml`

3. **Admin UI** (`admin/src/views/config/module/router/options.vue`):
   - CORS Origins input field has been added (line 3-9)
   - Binds to `data.cors_origins`

4. **Nginx Config** (`deploy/nginx/api.l9.lc.conf`):
   - Does not override CORS headers
   - Correctly proxies to backend

### The Issue
The `cors_origins` configuration cannot be set through the admin panel because:
1. The field exists in the UI but may not be properly saved to database
2. The default value is empty, blocking all non-localhost origins
3. There's no way to configure this at deployment time

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure CORS Origins via Admin Panel (Priority: P1)

As a system administrator, I need to configure allowed CORS origins through the admin panel so that my frontend application can communicate with the backend API.

**Why this priority**: Without CORS configuration, the frontend cannot make API calls, making the entire system non-functional.

**Independent Test**: Can be tested by logging into admin panel, navigating to router settings, adding a CORS origin, saving, and verifying the setting persists.

**Acceptance Scenarios**:

1. **Given** I am logged into the admin panel, **When** I navigate to Settings → Router → Options, **Then** I see a "CORS Origins" input field
2. **Given** I enter "https://www.l9.lc" in the CORS Origins field, **When** I save the settings, **Then** the value is persisted to the database
3. **Given** CORS origins are configured, **When** the backend receives a preflight OPTIONS request with matching origin, **Then** it returns proper CORS headers

---

### User Story 2 - Verify CORS Headers in API Response (Priority: P2)

As a developer, I need to verify that CORS headers are correctly returned so that I can debug cross-origin issues.

**Why this priority**: Essential for troubleshooting and confirming the fix works.

**Independent Test**: Can be tested by making a curl request with Origin header and inspecting response headers.

**Acceptance Scenarios**:

1. **Given** CORS origins include "https://www.l9.lc", **When** I send OPTIONS request with Origin: https://www.l9.lc, **Then** response includes `Access-Control-Allow-Origin: https://www.l9.lc`
2. **Given** CORS origins include "https://www.l9.lc", **When** I send OPTIONS request with Origin: https://evil.com, **Then** response does NOT include `Access-Control-Allow-Origin` header

---

### User Story 3 - Frontend Login Works After CORS Fix (Priority: P3)

As an end user, I need to log into the application from the frontend so that I can access protected features.

**Why this priority**: This is the ultimate validation that the entire system works end-to-end.

**Independent Test**: Can be tested by visiting https://www.l9.lc, entering credentials, and successfully logging in.

**Acceptance Scenarios**:

1. **Given** CORS is properly configured, **When** I visit https://www.l9.lc and click login, **Then** the login request succeeds without CORS errors
2. **Given** I am logged in, **When** I refresh the page, **Then** my session persists

---

### Edge Cases

- What happens when CORS origins field is empty? (Should default to localhost only for development)
- What happens when multiple origins are configured? (Should be comma-separated)
- What happens when origin has trailing slash? (Should handle gracefully)
- What happens with wildcard origin `*`? (Should work for public APIs, but not with credentials)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow administrators to configure CORS allowed origins through the admin panel
- **FR-002**: System MUST persist CORS origins configuration to the database
- **FR-003**: System MUST return proper CORS headers (`Access-Control-Allow-Origin`, `Access-Control-Allow-Credentials`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`) for allowed origins
- **FR-004**: System MUST support multiple comma-separated origins
- **FR-005**: System MUST handle preflight OPTIONS requests correctly
- **FR-006**: System MUST allow localhost origins by default for development
- **FR-007**: System MUST NOT expose CORS headers for disallowed origins

### Key Entities

- **Router Config**: Contains `cors_origins` field (comma-separated list of allowed origins)
- **CORS Middleware**: Intercepts requests and adds appropriate headers based on origin matching

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Administrators can configure CORS origins through admin panel in under 1 minute
- **SC-002**: API returns correct CORS headers for 100% of requests from allowed origins
- **SC-003**: Frontend login succeeds without CORS errors
- **SC-004**: System correctly rejects requests from non-allowed origins (no CORS headers returned)

## Assumptions

- Admin panel authentication is working correctly
- Database connection is stable
- Nginx configuration does not interfere with CORS headers
- Cloudflare CDN passes through CORS headers without modification
- Users have administrative access to configure settings
