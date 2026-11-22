## ADDED Requirements

### Requirement: Live log streaming
The panel and CLI SHALL stream a service's stdout and stderr in real time, across all its replicas, labelled by replica.

#### Scenario: Follow logs
- **WHEN** the admin opens a service's logs while it is serving traffic
- **THEN** new lines appear within one second of being written

### Requirement: Log retention
Service logs SHALL be stored compressed per service and per day and kept for 7 days by default (configurable).

#### Scenario: Investigating last night
- **WHEN** a service crashed 10 hours ago
- **THEN** the admin can read the logs from before the crash

#### Scenario: Retention expiry
- **WHEN** logs are older than the retention period
- **THEN** they are deleted automatically

### Requirement: Log search
Stored logs SHALL be searchable by text and time range, and by level when the app writes JSON lines with a level field.

#### Scenario: Filtered search
- **WHEN** the admin searches "timeout" between 03:00 and 04:00 with level "error"
- **THEN** only matching lines from that window are returned

### Requirement: Build logs
Each deploy SHALL keep its build log for as long as the deploy's image is kept.

#### Scenario: Image removed
- **WHEN** a deploy's image is removed by cleanup
- **THEN** its build log is removed with it
