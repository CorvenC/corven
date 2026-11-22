## ADDED Requirements

### Requirement: Single public API
The panel and the CLI SHALL use only the public HTTP API, which SHALL be documented by an OpenAPI description served by the platform.

#### Scenario: OpenAPI document
- **WHEN** a client requests the OpenAPI document
- **THEN** it describes every endpoint the panel and CLI use

### Requirement: Versioning
The API SHALL be served under `/api/v1`. Breaking changes MAY happen during v0.x; from v1.0 they SHALL require a new version prefix.

#### Scenario: Versioned path
- **WHEN** a client calls `/api/v1/projects`
- **THEN** it receives the project list

### Requirement: Scoped API tokens
API tokens SHALL have a scope (read-only, deploy, or full; optionally limited to one project) and an expiry date, and SHALL be shown only once at creation.

#### Scenario: Read-only token
- **WHEN** a read-only token calls a deploy endpoint
- **THEN** the API responds 403

#### Scenario: Expired token
- **WHEN** a token past its expiry date is used
- **THEN** the API responds 401

### Requirement: Error format
Errors SHALL use the RFC 9457 problem details format, including per-field validation errors.

#### Scenario: Validation error
- **WHEN** a request has an invalid field
- **THEN** the response is `application/problem+json` naming the field and the reason

### Requirement: Identifiers
Resource identifiers SHALL be type-prefixed and time-ordered (for example `prj_...`, `svc_...`, `dep_...`).

#### Scenario: Recognizable ID
- **WHEN** a deploy is created
- **THEN** its ID starts with `dep_`

### Requirement: Cursor pagination
List endpoints SHALL paginate with opaque cursors.

#### Scenario: Next page
- **WHEN** a client requests a list with the cursor from the previous response
- **THEN** it receives the following items without duplicates or gaps, even if new items were created meanwhile

### Requirement: Rate limiting
The API SHALL limit requests per token and per IP, returning 429 with a retry hint when exceeded.

#### Scenario: Runaway script
- **WHEN** a client exceeds the limit
- **THEN** it receives 429 with a `Retry-After` header

### Requirement: Real-time events
Live data (logs, build progress, deploy status, metrics) SHALL be delivered to clients as Server-Sent Events.

#### Scenario: Deploy progress
- **WHEN** a client subscribes to a deploy's events
- **THEN** it receives each step change as it happens
