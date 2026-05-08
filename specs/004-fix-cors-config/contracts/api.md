# API Contract: CORS Configuration

**Date**: 2026-05-08
**Feature**: 004-fix-cors-config

## Admin API: Router Config

### GET /admin/api/config/router

Retrieves the router configuration including CORS origins.

**Response**:
```json
{
  "admin_path": "/admin",
  "sitemap_path": "/sitemap",
  "cors_origins": "https://www.l9.lc",
  ...
}
```

### POST /admin/api/config/router

Updates the router configuration.

**Request Body**:
```json
{
  "cors_origins": "https://www.l9.lc,https://app.l9.lc",
  ...
}
```

**Response**:
```json
{
  "success": true
}
```

## Public API: CORS Headers

### OPTIONS /api/* (Preflight Request)

Handles CORS preflight requests.

**Request Headers**:
- `Origin: https://www.l9.lc`
- `Access-Control-Request-Method: POST`
- `Access-Control-Request-Headers: content-type`

**Response Headers** (if origin is allowed):
```
Access-Control-Allow-Origin: https://www.l9.lc
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization, X-Requested-With
Access-Control-Max-Age: 86400
Access-Control-Expose-Headers: Set-Cookie
Vary: Origin, Access-Control-Request-Method, Access-Control-Request-Headers
```

**Response Status**: `204 No Content`

### Actual Request (POST, GET, etc.)

**Response Headers** (if origin is allowed):
```
Access-Control-Allow-Origin: https://www.l9.lc
Access-Control-Allow-Credentials: true
Access-Control-Expose-Headers: Set-Cookie
Vary: Origin
```

## Error Cases

| Scenario | Response |
|----------|----------|
| Origin not in allowed list | No CORS headers (browser blocks request) |
| No Origin header | Request proceeds normally |
| Empty cors_origins config | Only localhost origins allowed |