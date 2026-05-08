# Data Model: Fix UI Theme Toggle and User Registration

**Feature**: Fix UI Theme Toggle and User Registration
**Date**: 2026-05-08

## Entities

### User Entity (Existing - Needs Migration Fix)

**Location**: `main/domain/core/entity/user.go`

**Fields**:

| Field | Type | GORM Tags | Description |
|-------|------|-----------|-------------|
| ID | uint | primaryKey, autoIncrement | Primary key |
| Username | string | uniqueIndex, not null, varchar(100) | Unique username |
| Email | string | uniqueIndex, not null, varchar(150) | Unique email address |
| Password | string | not null, varchar(250) | Bcrypt hashed password |
| Role | string | default:'user', varchar(20) | User role (user/admin) |
| CreatedAt | time.Time | default:CURRENT_TIMESTAMP | Record creation time |
| UpdatedAt | time.Time | default:CURRENT_TIMESTAMP | Record update time |

**Table Name**: `user`

**Indexes**:
- Primary key on `id`
- Unique index on `username`
- Unique index on `email`

**Relationships**: None (standalone entity for authentication)

**State Transitions**: None (simple CRUD entity)

---

### Theme State (Frontend - No Database Entity)

**Location**: `admin/src/store/index.js`

**Fields**:

| Field | Type | Storage | Description |
|-------|------|---------|-------------|
| dark | boolean | localStorage | Dark mode enabled |
| color | string | localStorage | Accent color (hex) |
| bgColor | string | localStorage | Light mode background color |
| darkBgColor | string | localStorage | Dark mode background color |

**Persistence**: Via `@vueuse/core` `useStorage` - stored in browser localStorage

---

## Validation Rules

### User Registration

- **Username**: Required, non-empty, unique
- **Email**: Required, non-empty, unique, valid email format
- **Password**: Required, minimum 6 characters

### Theme State

- **dark**: Boolean, defaults to system preference
- **color**: Hex color string (e.g., "#4CA1F7")
- **bgColor**: Hex color string
- **darkBgColor**: Hex color string

---

## Migration Requirements

### User Table Migration

**Status**: NOT IMPLEMENTED (Root Cause of Registration Failure)

**Required Changes**:

1. Add `MigrateTable()` method to `UserRepository`:
   ```go
   func (r *UserRepository) MigrateTable() error {
       return db.DB.AutoMigrate(&entity.User{})
   }
   ```

2. Add User migration call to `repository.MigrateTable()` in `main/domain/core/repository/repository.go`:
   ```go
   if err := User.MigrateTable(); err != nil {
       log.Error("migrate user table error", log.Err(err))
   }
   ```

**Migration Order**: After database initialization, before any user operations

---

## Data Flow

### User Registration Flow

```
Frontend (POST /api/auth/register)
    ↓
Controller (auth.go)
    ↓
Service (user.go: Register)
    ↓
Repository (user.go: Create)
    ↓
Database (user table)
```

### Theme Toggle Flow

```
User clicks toggle (Dark.vue)
    ↓
store.dark updated (Pinia)
    ↓
watch() triggered (base.vue)
    ↓
document.body.setAttribute/removeAttribute('arco-theme', 'dark')
    ↓
CSS applies based on arco-theme attribute
```
