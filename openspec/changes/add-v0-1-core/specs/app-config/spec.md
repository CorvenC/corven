## ADDED Requirements

### Requirement: Encrypted variables
Environment variables SHALL be stored encrypted at rest with a key derived from the master key, and secret values SHALL never be shown in full after being saved.

#### Scenario: Viewing a secret
- **WHEN** the admin opens a service's variables
- **THEN** secret values are masked, and revealing one requires re-authentication

#### Scenario: Database dump
- **WHEN** someone reads the internal database directly without the master key
- **THEN** variable values are unreadable

### Requirement: Shared variables
Variables MAY be defined at environment level and SHALL apply to every service in that environment unless the service defines the same name.

#### Scenario: Override
- **WHEN** the environment defines `LOG_LEVEL=info` and service `web` defines `LOG_LEVEL=debug`
- **THEN** `web` receives `debug` and other services receive `info`

### Requirement: References between services
A variable value MAY reference another service's variable using `${{service.NAME}}`, and the platform SHALL resolve references when the service starts.

#### Scenario: Reference updates
- **WHEN** the referenced value changes and the referencing service is restarted
- **THEN** it receives the new value without being edited

#### Scenario: Broken reference
- **WHEN** a reference points to a service or variable that does not exist
- **THEN** saving is rejected with an error naming the broken reference

### Requirement: Bulk import
The panel and CLI SHALL accept pasted `.env` content and parse it into individual variables.

#### Scenario: Paste .env
- **WHEN** the admin pastes five `KEY=value` lines
- **THEN** five variables are created, and invalid lines are reported without discarding the valid ones

### Requirement: Pending changes
Variable changes SHALL stay pending until applied; applying SHALL restart the service with zero downtime and without rebuilding.

#### Scenario: Several edits
- **WHEN** the admin edits three variables and clicks apply
- **THEN** the service restarts once with all three changes

### Requirement: Secret masking in logs
Values of variables marked secret SHALL be replaced with `****` in stored and streamed logs.

#### Scenario: App prints a secret
- **WHEN** an app writes the value of `STRIPE_KEY` to stdout
- **THEN** the log shows `****` instead of the value
