# API Contracts: Fix UI Theme Toggle and User Registration

**Feature**: Fix UI Theme Toggle and User Registration
**Date**: 2026-05-08

## Registration API

### POST /api/auth/register

**Description**: Register a new user account

**Request**:

```json
{
  "username": "string (required)",
  "email": "string (required, valid email)",
  "password": "string (required, min 6 characters)"
}
```

**Response (Success)**:

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

**Response (Error - Validation)**:

```json
{
  "error": "username, email, and password are required"
}
```

```json
{
  "error": "password must be at least 6 characters"
}
```

**Response (Error - Duplicate)**:

```json
{
  "success": false,
  "error": "username already exists"
}
```

```json
{
  "success": false,
  "error": "email already exists"
}
```

**Response (Error - Database)**:

```json
{
  "success": false,
  "error": "failed to create user"
}
```

**HTTP Status Codes**:
- 200: Success
- 400: Validation error or duplicate user
- 500: Server error (token generation failed)

---

## Theme API (No Backend API)

Theme is managed entirely on the frontend via:
- Pinia store (`admin/src/store/index.js`)
- localStorage persistence
- CSS attribute switching (`arco-theme` on body element)

No backend API changes required for theme functionality.
