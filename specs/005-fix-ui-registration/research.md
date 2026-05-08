# Research: Fix UI Theme Toggle and User Registration

**Feature**: Fix UI Theme Toggle and User Registration
**Date**: 2026-05-08
**Status**: Complete

## Research Questions

### RQ1: Why does dark/light mode toggle not display correctly?

**Context**: The admin panel uses Arco Design Vue with dark mode support via `arco-theme="dark"` attribute on the body element.

**Investigation**:

1. **Current Implementation** (`admin/src/layout/base.vue:32-38`):
   ```javascript
   // 初始化 dark
   if(store.dark) document.body.setAttribute('arco-theme', 'dark');
   // 监听 dark
   watch(()=>store.dark,(val)=>{
     if(val) document.body.setAttribute('arco-theme', 'dark');
     else document.body.removeAttribute('arco-theme');
   })
   ```

2. **Store State** (`admin/src/store/index.js:16`):
   ```javascript
   dark: useStorage("dark", useNavigatorDark()),
   ```

3. **Dark Mode Toggle** (`admin/src/components/app/Dark.vue`):
   ```vue
   <a-switch type="line" v-model="store.dark">
   ```

**Decision**: The implementation logic is correct. The issue may be:
- CSS styles not applying correctly for light mode
- Some components not responding to theme change
- Need to compare with original code in `xxc.zip`

**Action**: Extract and compare `xxc.zip` to identify differences in theme implementation.

---

### RQ2: Why does user registration fail with "failed to create user"?

**Context**: Registration endpoint returns `{"error":"failed to create user","success":false}`.

**Investigation**:

1. **Error Source** (`main/domain/core/service/user.go:54-57`):
   ```go
   err = repository.User.Create(user)
   if err != nil {
       return nil, errors.New("failed to create user")
   }
   ```

2. **Repository Create** (`main/domain/core/repository/user.go:14-16`):
   ```go
   func (r *UserRepository) Create(user *entity.User) error {
       return db.DB.Create(user).Error
   }
   ```

3. **User Entity** (`main/domain/core/entity/user.go:6-14`):
   ```go
   type User struct {
       ID        uint      `gorm:"type:int;size:32;primaryKey;autoIncrement"`
       Username  string    `gorm:"type:varchar(100);uniqueIndex;not null"`
       Email     string    `gorm:"type:varchar(150);uniqueIndex;not null"`
       Password  string    `gorm:"type:varchar(250);not null"`
       Role      string    `gorm:"type:varchar(20);default:'user'"`
       CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
       UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
   }
   ```

**Possible Causes**:
1. **Table doesn't exist**: GORM auto-migrate may not have run for User entity
2. **Database connection issue**: DB connection not initialized properly
3. **Migration not registered**: User entity not added to auto-migrate list

**Decision**: Check if User entity is registered in the startup/auto-migrate process.

**Action**: Search for where GORM auto-migrate is called and verify User entity is included.

---

## Findings

### Finding 1: Theme Implementation Appears Correct

The theme toggle implementation follows Arco Design's recommended approach:
- `arco-theme="dark"` attribute on body for dark mode
- Remove attribute for light mode
- State persisted in localStorage via Pinia + vueuse

**Next Step**: Compare with `xxc.zip` to identify any missing styles or component configurations.

### Finding 2: User Entity Migration NOT Registered (ROOT CAUSE FOUND)

**CONFIRMED**: User entity is NOT registered for auto-migration!

The `main/domain/core/repository/repository.go` file contains `MigrateTable()` function that calls `MigrateTable()` for:
- Article
- Category
- Tag
- Mapping
- Link
- Store

**But User is MISSING!**

Additionally, `main/domain/core/repository/user.go` does NOT have a `MigrateTable()` method, unlike other repositories.

**Root Cause**: The `user` table is never created because:
1. `UserRepository` has no `MigrateTable()` method
2. `repository.MigrateTable()` doesn't call User migration

**Fix Required**:
1. Add `MigrateTable()` method to `UserRepository`
2. Add `User.MigrateTable()` call to `repository.MigrateTable()`

---

## Technical Decisions

### TD1: Theme Fix Approach

**Decision**: Compare current implementation with `xxc.zip` reference code
**Rationale**: User provided reference code, likely contains working implementation
**Alternatives Considered**:
- Debug CSS manually - time consuming, may miss edge cases
- Rebuild from scratch - unnecessary if only minor differences

### TD2: Registration Fix Approach

**Decision**: Verify database migration and add User entity to auto-migrate if missing
**Rationale**: Most likely cause is missing table or migration
**Alternatives Considered**:
- Add better error logging - helpful but doesn't fix root cause
- Manual table creation - not maintainable, should use GORM auto-migrate

---

## Next Steps

1. Extract `xxc.zip` and compare theme-related files
2. ~~Check `main/startup/startup.go` for User entity auto-migration~~ **DONE - ROOT CAUSE FOUND**
3. Add `MigrateTable()` method to `UserRepository`
4. Add `User.MigrateTable()` call to `repository.MigrateTable()`
5. Test database connection and table existence
