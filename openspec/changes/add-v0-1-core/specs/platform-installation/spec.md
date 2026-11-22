## ADDED Requirements

### Requirement: One-command installation
The platform SHALL install on a supported server with a single shell command that installs Docker if missing and starts the platform, its database and its reverse proxy as containers.

#### Scenario: Fresh Ubuntu server
- **WHEN** an operator runs the install command on a clean Ubuntu LTS or Debian stable server (amd64 or arm64)
- **THEN** the platform, its internal PostgreSQL and Caddy are running as containers
- **AND** the installer prints the URL of the setup wizard

#### Scenario: Unsupported system
- **WHEN** the install command runs on an unsupported OS or CPU architecture
- **THEN** the installer stops before changing anything and explains which systems are supported

### Requirement: Release verification
The installer SHALL verify the signature of the platform release before installing it.

#### Scenario: Tampered artifact
- **WHEN** the downloaded release does not match its published signature
- **THEN** the installer aborts and changes nothing on the server

### Requirement: Resource footprint
The platform's own processes (excluding user apps) SHALL use less than 300 MB of RAM when idle, so it fits a server with 2 GB of RAM and 2 CPUs.

#### Scenario: Idle platform
- **WHEN** the platform is running with no deploys in progress
- **THEN** the combined memory of the platform, its database and Caddy stays under 300 MB

### Requirement: Firewall setup
The installer SHALL offer to restrict inbound traffic to ports 22, 80 and 443 and SHALL prevent Docker-published ports from bypassing that firewall.

#### Scenario: Operator accepts firewall setup
- **WHEN** the operator accepts the firewall prompt
- **THEN** only ports 22, 80 and 443 accept inbound connections, including for ports published by containers

#### Scenario: Operator declines
- **WHEN** the operator declines the firewall prompt
- **THEN** the installer leaves firewall rules unchanged

### Requirement: Master key generation
The installer SHALL generate a master encryption key, store it on the server outside the database, and show it once so the operator can keep a copy.

#### Scenario: Key shown once
- **WHEN** installation completes
- **THEN** the master key is displayed once with a warning that encrypted data cannot be recovered on another server without it

### Requirement: Setup wizard
On first access the panel SHALL guide the operator through creating the admin account, configuring a domain (or the no-domain mode), connecting GitHub or GitLab, and deploying a sample app.

#### Scenario: First visit
- **WHEN** the operator opens the panel before any admin exists
- **THEN** the setup wizard starts and no other page is reachable until the admin is created

### Requirement: Clean uninstall
The platform SHALL provide an uninstall command that removes the platform and asks whether to also remove user apps and their data.

#### Scenario: Keep apps
- **WHEN** the operator uninstalls and chooses to keep apps
- **THEN** the platform containers are removed and user app containers and volumes keep running
