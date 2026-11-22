## ADDED Requirements

### Requirement: Password login
The panel SHALL authenticate the admin with a username and password, storing only an Argon2 hash of the password.

#### Scenario: Valid credentials
- **WHEN** the admin submits the correct username and password
- **THEN** the panel asks for the 2FA code if 2FA is enabled, or starts a session otherwise

#### Scenario: Invalid credentials
- **WHEN** the submitted credentials are wrong
- **THEN** login fails with a generic message that does not reveal whether the username exists

### Requirement: TOTP two-factor authentication
The panel SHALL support TOTP 2FA and SHALL issue 10 single-use recovery codes when 2FA is enabled.

#### Scenario: Login with 2FA
- **WHEN** 2FA is enabled and the admin enters a valid TOTP code or an unused recovery code
- **THEN** a session starts, and a used recovery code can never be used again

### Requirement: Brute-force protection
The panel SHALL apply an increasing delay after 5 failed login attempts from the same IP or for the same account.

#### Scenario: Repeated failures
- **WHEN** a sixth failed attempt arrives within the window
- **THEN** the request is rejected until the delay expires, and the delay grows with further failures

### Requirement: Sessions
Sessions SHALL use secure, httpOnly cookies, last at most 7 days, expire after 12 hours of inactivity, and be listable and revocable by the admin.

#### Scenario: Inactive session
- **WHEN** a session has not been used for 12 hours
- **THEN** the next request requires logging in again

#### Scenario: Revoke another device
- **WHEN** the admin revokes a session from the sessions list
- **THEN** that session's next request is rejected

### Requirement: Access recovery from the server
An operator with shell access to the server SHALL be able to generate a one-time password reset link with a CLI command, without email.

#### Scenario: Forgotten password
- **WHEN** the operator runs the reset command on the server
- **THEN** it prints a single-use reset link that expires in 15 minutes

### Requirement: Dangerous action confirmation
Destructive or security-critical actions SHALL require typing the resource name and, when 2FA is enabled, a fresh TOTP code.

#### Scenario: Deleting a service
- **WHEN** the admin deletes a service
- **THEN** the panel requires the exact service name and a TOTP code before proceeding

### Requirement: Team-ready ownership model
Every resource SHALL have an owner and the data model SHALL support roles, even though v0.1 exposes a single admin.

#### Scenario: Single admin in v0.1
- **WHEN** the admin creates a project
- **THEN** the project is recorded as owned by the admin's account
