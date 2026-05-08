# Feature Specification: Fix UI Theme Toggle and User Registration

**Feature Branch**: `005-fix-ui-registration`
**Created**: 2026-05-08
**Status**: Draft
**Input**: User description: "问题1. UI还原度，当前暗黑模式，切换成亮色模式时，不正确， 如果需要参考，可以从 @xxc.zip 获取UI原始代码代码， 问题2，部署后，注册账号不成功，提示 {"error":"failed to create user","success":false}"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Dark/Light Mode Toggle Works Correctly (Priority: P1)

As a user of the admin panel, I want to switch between dark mode and light mode so that the UI displays correctly in both themes.

**Why this priority**: This is a visual bug that affects user experience and may indicate deeper issues with theme state management.

**Independent Test**: Can be fully tested by toggling the dark/light mode switch in the admin panel and verifying all UI elements display correctly in both modes.

**Acceptance Scenarios**:

1. **Given** the admin panel is in dark mode, **When** user toggles to light mode, **Then** all UI elements (backgrounds, text, cards, buttons) display with correct light theme colors
2. **Given** the admin panel is in light mode, **When** user toggles to dark mode, **Then** all UI elements display with correct dark theme colors
3. **Given** user refreshes the page, **When** the page loads, **Then** the theme state persists correctly (dark stays dark, light stays light)

---

### User Story 2 - User Registration Works After Deployment (Priority: P1)

As a new user, I want to register an account on the frontend so that I can access protected features.

**Why this priority**: Registration is a critical user flow that blocks new users from accessing the application.

**Independent Test**: Can be fully tested by visiting the frontend registration page, entering valid credentials, and verifying successful account creation.

**Acceptance Scenarios**:

1. **Given** a new user visits the registration page, **When** they enter valid username, email, and password (6+ characters), **Then** the account is created successfully
2. **Given** a user tries to register with an existing username, **When** they submit the form, **Then** an appropriate error message is displayed
3. **Given** a user tries to register with an existing email, **When** they submit the form, **Then** an appropriate error message is displayed
4. **Given** registration succeeds, **When** the response returns, **Then** the user is logged in automatically with a valid session

---

### Edge Cases

- What happens when the theme toggle is clicked rapidly multiple times?
- What happens when the browser's preferred color scheme changes while the page is open?
- What happens when registration fails due to database connection issues?
- What happens when password contains special characters?

## Requirements *(mandatory)*

### Functional Requirements

#### Theme Toggle (User Story 1)

- **FR-001**: System MUST apply `arco-theme="dark"` attribute to body element when dark mode is enabled
- **FR-002**: System MUST remove `arco-theme="dark"` attribute from body element when light mode is enabled
- **FR-003**: System MUST persist theme preference in localStorage via `useStorage("dark", ...)`
- **FR-004**: System MUST apply correct background colors for both themes (dark: `darkBgColor`, light: `bgColor`)
- **FR-005**: All Arco Design components MUST render correctly in both themes

#### User Registration (User Story 2)

- **FR-006**: System MUST accept username, email, and password for registration
- **FR-007**: System MUST validate that username, email, and password are non-empty
- **FR-008**: System MUST validate password length is at least 6 characters
- **FR-009**: System MUST check for existing username before creating user
- **FR-010**: System MUST check for existing email before creating user
- **FR-011**: System MUST hash password using bcrypt before storage
- **FR-012**: System MUST create user record in database with role "user"
- **FR-013**: System MUST return appropriate error message when registration fails

### Key Entities

- **User**: Represents a registered user with username, email, hashed password, role, and timestamps
- **Theme State**: Stores dark mode preference, background colors, and accent color

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Theme toggle switches visual appearance within 100ms of click
- **SC-002**: All UI elements display correctly in both themes with no visual artifacts
- **SC-003**: User can successfully register a new account in under 30 seconds
- **SC-004**: Registration returns clear error messages for validation failures
- **SC-005**: Theme preference persists across browser sessions

## Assumptions

- The original UI code in `xxc.zip` contains the correct theme implementation to reference
- The registration error "failed to create user" is returned from `service.User.Register` in `main/domain/core/service/user.go`
- The database connection is working (other operations like login work)
- The user table migration has run successfully (table exists with correct schema)
- Arco Design Vue components support both light and dark themes via `arco-theme` attribute
