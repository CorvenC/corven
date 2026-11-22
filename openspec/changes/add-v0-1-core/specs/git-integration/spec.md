## ADDED Requirements

### Requirement: GitHub connection
The platform SHALL connect to GitHub through a GitHub App, so repository access is granted per repository and webhooks are registered automatically.

#### Scenario: Connect a GitHub repo
- **WHEN** the admin installs the GitHub App and selects a repository
- **THEN** the repository appears as a deploy source and pushes to it are received as webhooks

### Requirement: GitLab connection
The platform SHALL connect to GitLab (gitlab.com or self-managed) through OAuth and register project webhooks automatically.

#### Scenario: Connect a GitLab repo
- **WHEN** the admin authorizes GitLab and selects a project
- **THEN** the project appears as a deploy source and a webhook is created on it

### Requirement: Deploy keys fallback
The platform SHALL support any Git URL accessed with a per-repository SSH deploy key that it generates.

#### Scenario: Generic Git server
- **WHEN** the admin adds a Git URL with a deploy key
- **THEN** the platform shows the public key to add on the Git server and can clone once it is added

### Requirement: Webhook authenticity
The platform SHALL reject webhooks whose signature or secret token does not verify.

#### Scenario: Forged webhook
- **WHEN** a request arrives at the webhook endpoint with an invalid signature
- **THEN** it is rejected with 401 and no deploy is triggered

### Requirement: Commit status reporting
The platform SHALL report build and deploy status on the commit in GitHub or GitLab, linking to the deploy page.

#### Scenario: Successful deploy
- **WHEN** a deploy triggered by a commit finishes successfully
- **THEN** the commit shows a success status linking to the deploy

#### Scenario: Failed deploy
- **WHEN** the build or rollout fails
- **THEN** the commit shows a failure status linking to the deploy logs
