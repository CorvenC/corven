# Design

## Context

A self-hosted, open source (AGPL-3.0) deployment platform for solo developers and small teams with one VPS. v0.1 runs on a single Ubuntu/Debian server (amd64/arm64) with at least 2 GB of RAM and 2 CPUs, and must stay under 300 MB itself. It is built by one developer with AI assistance, in variable weekly hours, with a target of 2026-12-31. Key decisions are recorded as ADRs in `docs/adr/`.

## Goals / Non-Goals

**Goals:**
- Push-to-deploy loop end to end: install, connect repo, push, HTTPS app with logs and rollback.
- Security above the norm for this category: socket proxy, encrypted secrets, signed releases, 2FA.
- Light footprint and a Railway-like experience (project canvas, service references).
- English and Spanish UI from day one.

**Non-Goals:**
- Multi-server, teams and roles UI, previews, managed databases, backups, templates, compose import, workers/cron, notifications, audit log (all v1, after v0.1).
- Kubernetes, serverless functions, Windows servers (never).
- Telemetry of any kind.

## Decisions

### Architecture: modular monolith, everything in containers
One Go binary with internal modules (`auth`, `git`, `build`, `deploy`, `proxy`, `logs`, `metrics`, `api`). It runs as a container next to its PostgreSQL and Caddy containers. Alternatives: separate agent per server (premature until multi-server), microservices (too much for one developer). Module boundaries are kept clean so an agent can be split out later. See ADR-0001.

### Repository layout
```
cmd/server        platform binary
cmd/cli           CLI binary
internal/<module> one package per module
web/              React + Vite panel (embedded with go:embed)
docs/             Starlight site (en, es) + docs/adr
deploy/           installer script, compose files
examples/         sample app for the setup wizard
```

### Backend stack
- HTTP and API: chi router + Huma, which generates the OpenAPI document and request validation from Go structs (ADR-0003).
- Database: PostgreSQL in its own container, accessed with sqlc + pgx; migrations with goose, embedded in the binary and applied at startup.
- Jobs: River on PostgreSQL for builds, deploys, cleanup and future cron, so jobs survive restarts.
- Logging: `log/slog` with JSON output in production.
- Docker: official Go SDK, talking to the daemon through a socket proxy that allows only the endpoints the platform needs (ADR-0004).
- Builds: BuildKit with layer cache; images stay in the local Docker image store (no registry in v0.1).

### Proxy and TLS
Caddy runs in its own container and is configured only through its admin API on an internal network (ADR-0002). The platform pushes the full route config on every change and re-pushes it after Caddy restarts. Wildcard certificates use the DNS-01 challenge with the Cloudflare module; sslip.io hostnames use HTTP-01. Custom domains get a certificate only after the platform verifies DNS.

### Zero-downtime rollout
1. Start new containers (all replicas) on the environment network with a new name suffix.
2. Run the health check against each replica until it passes or the grace period (default 60 s) ends.
3. Update the Caddy upstreams to the new replicas.
4. Wait for in-flight requests to drain (default 10 s), then stop the old containers.
On failure at step 2, remove the new containers and leave routes unchanged.

### Desired state and reconciliation
The database is the source of truth for what should run. A reconciler compares it with Docker every 30 s and on startup, creating, restarting or removing containers as needed, in dependency order. Crash handling (backoff, crash-loop state) lives in the reconciler, not in Docker restart policies.

### Secrets
AES-256-GCM with per-value nonces. A data key is encrypted by the master key, which lives in a file on the host, outside the database and backups. Rotating the master key re-encrypts the data key only. See ADR-0005.

### Logs
A collector follows container logs through the Docker API, masks known secret values, fans out to live SSE subscribers, and appends to per-service, per-day files compressed with zstd, plus a small time index for range queries. Search scans the files in the requested window.

### Metrics
Sampled every 15 s from the Docker stats API and the host (`/proc`), stored in PostgreSQL tables: raw for 24 h, hourly aggregates for 30 days. HTTP metrics come from Caddy access logs.

### Frontend
React + Vite + TypeScript (strict), Tailwind + shadcn/ui, TanStack Router and TanStack Query, TypeScript client generated from the OpenAPI document, Lingui for i18n (en, es), uPlot for charts. Dark, sober visual style with a light mode. Performance budget: initial load under 1 s, navigation under 100 ms. Accessibility target: WCAG AA.

### Auth
Argon2id password hashes, TOTP (RFC 6238) with 10 recovery codes, server-side sessions in PostgreSQL referenced by a secure httpOnly cookie. CLI login uses an OAuth-style device flow against the platform itself.

### Quality and supply chain
Trunk-based development with small MRs (about 400 lines max), Conventional Commits with scope, DCO sign-off. CI: golangci-lint, Go unit tests and integration tests with real Docker (testcontainers), Vitest, Playwright on key flows, govulncheck, gitleaks, gosec, CodeQL, Trivy. Releases: SemVer, changelog from commits, GoReleaser, cosign signatures and SBOM; images published to GHCR and the GitLab registry.

## Risks / Trade-offs

- **Scope versus deadline** → v0.1 is cut to the walking skeleton plus 2FA, encrypted variables, basic metrics and CLI; everything else moves to v1. Tasks are ordered so the project is usable at every checkpoint.
- **Docker socket is root-equivalent** → socket proxy with an allowlist; the platform container itself runs unprivileged.
- **Docker bypasses ufw** → the installer adds rules to the `DOCKER-USER` chain.
- **Postgres adds RAM** → tuned for small servers (shared_buffers about 64 MB); measured in CI against the 300 MB budget.
- **Let's Encrypt rate limits** → wildcard certificate for automatic subdomains; custom domains only after DNS verification.
- **Developer is learning Go** → small MRs, every one read and approved by the developer; ADRs explain the why.

## Migration Plan

Not applicable: first version. From v0.1 on, upgrades back up the internal database, pull the new signed image, apply embedded migrations and roll back automatically if the new version fails its health check.

## Open Questions

- Exact keys of the repository configuration file `corven.toml` (name decided: product is Corven).
- Whether the socket proxy is an existing project (for example Tecnativa's docker-socket-proxy) or a small built-in one.
