## ADDED Requirements

### Requirement: Supported systems
The CLI SHALL run on Windows, macOS and Linux (amd64 and arm64) and SHALL be installable by script (shell and PowerShell), Homebrew, Scoop/winget and `go install`.

#### Scenario: Windows install
- **WHEN** a user runs the PowerShell install script on Windows
- **THEN** the CLI is installed and available on the PATH

### Requirement: Browser login
`login` SHALL authenticate through the browser with a device code approved in the panel, and SHALL never ask for the password in the terminal. Non-interactive use SHALL accept an API token from an environment variable.

#### Scenario: Device login
- **WHEN** the user runs `login`
- **THEN** the CLI shows a code and opens the panel, and after approval it stores a token

#### Scenario: CI usage
- **WHEN** the token environment variable is set
- **THEN** commands authenticate with it without prompting

### Requirement: Core commands
The CLI SHALL provide `deploy`, `logs` (with follow), `env` (list, set, unset, import), and `rollback`.

#### Scenario: Deploy and follow
- **WHEN** the user runs `deploy` for a service
- **THEN** the CLI streams build progress and exits non-zero if the deploy fails

### Requirement: Machine-readable output
Every command SHALL support `--json`; without it, output SHALL be human-readable tables with color when the terminal supports it.

#### Scenario: Scripting
- **WHEN** the user runs a list command with `--json`
- **THEN** stdout contains only valid JSON
