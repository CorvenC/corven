# Proposal

## Why

Developers with a single VPS want the push-to-deploy comfort of Vercel or Railway without the bill or the lock-in, but the existing self-hosted options are heavy, weak on security, or clumsy to use. v0.1 proves the core loop end to end: install with one command, connect a repo, push, and the app is live over HTTPS with logs and rollback. Target date: 2026-12-31.

## What Changes

- One-command installer for Ubuntu/Debian (amd64 and arm64) that sets up Docker, the platform, its Postgres and Caddy, all as containers.
- A single admin account with password, TOTP 2FA, recovery codes and brute-force protection.
- GitHub App and GitLab OAuth integration: webhooks trigger deploys and the platform reports commit status.
- Deploys from any repo that has a Dockerfile (configurable folder and Dockerfile for monorepos) or from a prebuilt image.
- Zero-downtime deploys gated by health checks, rollback to any of the last N images, cancellation of stale builds, build timeouts.
- Automatic subdomains (wildcard DNS via Cloudflare, or sslip.io without a domain), custom domains with DNS verification, automatic TLS.
- Encrypted environment variables, shared per environment, with references between services.
- Live logs plus 7 days of searchable history.
- Basic metrics (CPU, RAM, disk, network) for the server and each app.
- Public REST API (OpenAPI, versioned under `/api/v1`) with scoped, expiring tokens.
- A basic CLI (login, deploy, logs, env, rollback) for Windows, macOS and Linux.
- Desired-state reconciliation so apps come back after crashes and reboots.

Out of scope for v0.1 (planned for v1): preview environments, managed Postgres/Redis, backups, app templates, docker-compose import, workers and cron, notifications, audit log, multi-user teams. Never in scope: Kubernetes, serverless functions, Windows servers.

## Capabilities

### New Capabilities
- `platform-installation`: installing, upgrading and uninstalling the platform on a server.
- `admin-auth`: admin account, login, 2FA, sessions, recovery and dangerous-action confirmation.
- `git-integration`: connecting GitHub and GitLab, receiving webhooks, reporting commit status.
- `app-deployment`: building images, rolling out deploys, health checks, rollback and runtime supervision.
- `routing-tls`: domains, subdomains, redirects, access restrictions and TLS certificates for apps.
- `app-config`: environment variables and secrets for services.
- `app-logs`: collecting, streaming, storing and searching app and build logs.
- `metrics`: collecting and displaying resource usage for the server and apps.
- `public-api`: the HTTP API, its tokens, versioning, errors and pagination.
- `cli`: the command-line client.

### Modified Capabilities

None (first change in the project).

## Impact

- New repository (GitLab, mirrored to GitHub), AGPL-3.0.
- Backend in Go (single modular binary), panel in React + Vite embedded in the binary.
- Runtime dependencies on the server: Docker Engine, Caddy (separate container), PostgreSQL (separate container), a Docker socket proxy.
- External services: GitHub, GitLab, Let's Encrypt, Cloudflare DNS API (optional), sslip.io (optional).
