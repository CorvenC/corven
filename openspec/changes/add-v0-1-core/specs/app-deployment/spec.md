## ADDED Requirements

### Requirement: Resource hierarchy
Resources SHALL be organized as Project > Environment > Service, where a service is an app built from a repo or an app run from a prebuilt image.

#### Scenario: New project
- **WHEN** the admin creates a project
- **THEN** it has a `production` environment, and more environments can be added

### Requirement: Dockerfile builds
The platform SHALL build a service from the Dockerfile in its repository, using a configurable build folder and Dockerfile path, with layer caching between builds.

#### Scenario: Monorepo service
- **WHEN** a service is configured with folder `apps/web` and Dockerfile `apps/web/Dockerfile`
- **THEN** the build uses that folder as context and that Dockerfile

#### Scenario: Unrelated change in a monorepo
- **WHEN** a push changes no file inside a service's build folder
- **THEN** that service is not redeployed automatically

#### Scenario: No Dockerfile
- **WHEN** the configured Dockerfile does not exist
- **THEN** the deploy fails with a message naming the detected project type and offering a tested Dockerfile template

### Requirement: Prebuilt image services
The platform SHALL deploy a service from a public or authenticated container image reference without a repository.

#### Scenario: Deploy an image
- **WHEN** the admin creates a service from image `louislam/uptime-kuma:1`
- **THEN** it runs with a subdomain, TLS, logs and metrics like any other service

### Requirement: Deploy triggers
A deploy SHALL be triggerable by a push to the service's configured branch, by a button in the panel, by the CLI, and by the API. Automatic deploys on push SHALL be disableable per service.

#### Scenario: Auto-deploy off
- **WHEN** auto-deploy is disabled and a push arrives
- **THEN** no deploy starts, and the panel shows the new commit as available to deploy

### Requirement: Build queue and limits
Builds SHALL run in a persistent queue with per-build CPU and memory limits and a timeout of 30 minutes by default (configurable).

#### Scenario: Build exceeds timeout
- **WHEN** a build runs longer than its timeout
- **THEN** it is stopped and the deploy is marked failed with reason "timeout"

#### Scenario: Platform restarts mid-build
- **WHEN** the platform restarts while a build is queued or running
- **THEN** the job is resumed or retried after the restart instead of being lost

### Requirement: Superseded builds
When a new deploy is requested for a service with a build in progress, the platform SHALL cancel the older build by default; services MAY be configured to queue instead.

#### Scenario: Two quick pushes
- **WHEN** commit B is pushed while commit A is still building
- **THEN** A's build is cancelled and B is built

### Requirement: Zero-downtime rollout
A new version SHALL receive traffic only after passing its health check; the previous version SHALL be stopped only after traffic has switched.

#### Scenario: Healthy new version
- **WHEN** the new containers pass the health check
- **THEN** traffic switches to them and the old containers stop, with no failed requests during the switch

#### Scenario: Unhealthy new version
- **WHEN** the new containers do not pass the health check within the grace period
- **THEN** they are removed, the old version keeps serving, and the deploy is marked failed

### Requirement: Health checks
The health check SHALL be an HTTP request to a configurable path (default `/health`) expecting a 2xx response, falling back to a TCP check of the service port when no path is configured.

#### Scenario: TCP fallback
- **WHEN** a service has no health path configured
- **THEN** it is considered healthy once its port accepts TCP connections

### Requirement: Port detection
The platform SHALL determine the service port from the configuration file or panel, else from the image's `EXPOSE`, and SHALL inject a `PORT` variable with that value.

#### Scenario: Port from EXPOSE
- **WHEN** no port is configured and the image exposes 3000
- **THEN** traffic is routed to port 3000 and `PORT=3000` is set

### Requirement: Repository configuration file
A service MAY be configured by an optional file in its repository (port, health check, resources, domains); values in the file SHALL override panel values, and the panel SHALL show which values come from the file.

#### Scenario: File present
- **WHEN** the repository contains the configuration file with `port = 8080`
- **THEN** the service uses port 8080 and the panel marks that setting as managed by the file

### Requirement: Rollback
The platform SHALL keep the images of the last 10 successful deploys per service (configurable) and SHALL roll back to any of them without rebuilding.

#### Scenario: Roll back
- **WHEN** the admin rolls back to a deploy from three versions ago
- **THEN** that image is rolled out with the same zero-downtime process, in seconds

### Requirement: Replicas
A service SHALL run 1 to N replicas on the server, with traffic balanced across healthy replicas.

#### Scenario: Scale to 3
- **WHEN** the admin sets replicas to 3
- **THEN** three containers run and requests are distributed among them

### Requirement: Resource limits
Each service SHALL run with memory and CPU limits, defaulting to 512 MB and 1 CPU, editable per service.

#### Scenario: Memory leak
- **WHEN** a service exceeds its memory limit
- **THEN** only that service's container is killed and restarted; other services are unaffected

### Requirement: Crash handling
A crashed service SHALL be restarted with exponential backoff; after repeated crashes in a short window it SHALL be marked failed and stop restarting.

#### Scenario: Crash loop
- **WHEN** a service crashes 5 times within 10 minutes
- **THEN** it is marked "crash loop" in the panel and no further automatic restarts happen until the admin acts

### Requirement: Isolation
Each environment SHALL have its own Docker network; services SHALL run unprivileged and SHALL NOT reach services of other projects over the network.

#### Scenario: Cross-project access
- **WHEN** a container in project A tries to connect to a service in project B by name or IP
- **THEN** the connection fails

### Requirement: Persistent volumes
A service MAY mount named volumes that survive deploys, rollbacks and restarts, and the panel SHALL show their size.

#### Scenario: Data survives deploy
- **WHEN** a service writes a file to its volume and is redeployed
- **THEN** the file is still present in the new version

### Requirement: Desired-state reconciliation
The platform SHALL store the desired state of every service and, at least every 30 seconds and after a server reboot, SHALL correct any difference between it and what is actually running.

#### Scenario: Container removed by hand
- **WHEN** someone deletes a service container directly with Docker
- **THEN** the platform recreates it within 30 seconds

#### Scenario: Server reboot
- **WHEN** the server reboots
- **THEN** all services come back, with services that others depend on started first

### Requirement: Disk cleanup
The platform SHALL automatically remove images not kept for rollback, stale build cache and dead containers, and SHALL warn when disk usage exceeds 85%.

#### Scenario: Disk nearly full
- **WHEN** disk usage passes 85%
- **THEN** the panel shows a warning on the home screen

### Requirement: Deploy page
Each deploy SHALL have a page with its steps and durations, the commits included since the previous deploy, the live build log, and, on failure, a short explanation of the probable cause.

#### Scenario: Missing variable
- **WHEN** a deploy fails because the app exits complaining about an unset variable
- **THEN** the deploy page names the missing variable as the probable cause

### Requirement: Soft delete
Deleted projects and services SHALL be stopped and kept recoverable, with their volumes, for 7 days before permanent removal.

#### Scenario: Restore within 7 days
- **WHEN** the admin restores a service deleted 3 days ago
- **THEN** it comes back with its configuration and volume data
