## ADDED Requirements

### Requirement: Resource metrics
The platform SHALL sample CPU, memory, disk and network usage of the server and of each service every 15 seconds.

#### Scenario: Service chart
- **WHEN** the admin opens a service's metrics
- **THEN** charts show its CPU, memory and network over time, per replica and in total

### Requirement: Metrics retention
Raw samples SHALL be kept for 24 hours and hourly aggregates for 30 days.

#### Scenario: Last month
- **WHEN** the admin views the last 30 days
- **THEN** hourly averages and peaks are shown

### Requirement: Platform health indicators
The panel SHALL show HTTP requests per minute, 5xx error rate and latency per service, restarts and crashes, build and deploy durations, and certificate expiry dates.

#### Scenario: Error spike
- **WHEN** a service starts returning 5xx responses
- **THEN** its error rate chart reflects it within one minute

### Requirement: Server overview
The home screen SHALL show the server's current CPU, memory and disk usage next to the project list with each project's status.

#### Scenario: Home screen
- **WHEN** the admin opens the panel
- **THEN** each project card shows healthy, deploying or down, and a strip shows server usage
