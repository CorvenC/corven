## 0. Before the first public commit

- [x] 0.1 Choose the final product name; check GitHub, GitLab and package-manager name availability (Corven)
- [x] 0.2 Define logo, colors and typography (docs/brand/BRAND.md)
- [x] 0.3 Install Ubuntu in WSL2 with Docker and Go (Ubuntu 24.04, native Docker Engine, Go, Task); test installs use a throwaway WSL distro instead of Multipass (needs Hyper-V or VirtualBox on Windows Home)
- [x] 0.4 Create the `corvenc` GitLab group and repo, the `CorvenC` GitHub org, pushing to both, AGPL-3.0 license, README, CONTRIBUTING, SECURITY, CODE_OF_CONDUCT

## 1. Walking skeleton (clone → build → run → domain)

- [x] 1.1 Go module, repo layout, Taskfile, golangci-lint config
- [x] 1.2 CI pipeline: lint, test, build for Go
- [x] 1.3 Docker client through the SDK: build an image from a local folder and run it
- [ ] 1.4 Clone a public Git repo at a given commit and build it with BuildKit
- [ ] 1.5 Caddy container with admin API; register a route to a running container
- [ ] 1.6 sslip.io hostname with automatic TLS for that route
- [ ] 1.7 Command `server deploy <git-url>` that does all of the above end to end (throwaway entry point)

## 2. Platform base

- [ ] 2.1 PostgreSQL container, pgx pool, goose migrations embedded
- [ ] 2.2 Data model: accounts, projects, environments, services, deploys (with owner and role fields)
- [ ] 2.3 sqlc queries for the model
- [ ] 2.4 River job queue with build and deploy jobs
- [ ] 2.5 Huma + chi API skeleton under `/api/v1`, OpenAPI served, RFC 9457 errors
- [ ] 2.6 Prefixed time-ordered IDs and cursor pagination helpers
- [ ] 2.7 slog JSON logging and config loading

## 3. Auth

- [ ] 3.1 Admin creation, Argon2id hashing, login and server-side sessions
- [ ] 3.2 Brute-force delays per IP and per account
- [ ] 3.3 TOTP enrollment, verification and recovery codes
- [ ] 3.4 Session list and revocation; inactivity expiry
- [ ] 3.5 Server-side password reset command
- [ ] 3.6 Scoped, expiring API tokens
- [ ] 3.7 Dangerous-action confirmation (name + TOTP)

## 4. Deploy engine

- [ ] 4.1 Build queue with per-build CPU/memory limits and timeout
- [ ] 4.2 Cancel superseded builds (queue mode as option)
- [ ] 4.3 Port detection (config, EXPOSE) and `PORT` injection
- [ ] 4.4 HTTP health check with TCP fallback
- [ ] 4.5 Zero-downtime rollout and failure handling
- [ ] 4.6 Replicas with balanced upstreams
- [ ] 4.7 Resource limits per service
- [ ] 4.8 Per-environment networks, unprivileged containers, cross-project isolation test
- [ ] 4.9 Named volumes and size reporting
- [ ] 4.10 Desired-state reconciler with dependency order and crash-loop detection
- [ ] 4.11 Rollback to any kept image; image retention (last 10)
- [ ] 4.12 Disk cleanup job and 85% warning
- [ ] 4.13 Prebuilt image services
- [ ] 4.14 Monorepo folder/Dockerfile settings and changed-path filter
- [ ] 4.15 Repository configuration file parsing and precedence
- [ ] 4.16 Dockerfile template suggestion when missing
- [ ] 4.17 Failure explanation rules (missing variable, port, memory, health check)
- [ ] 4.18 Soft delete with 7-day recovery

## 5. Git integration

- [ ] 5.1 GitHub App: install flow, repo listing, webhook verification
- [ ] 5.2 GitLab OAuth: project listing, webhook creation and verification
- [ ] 5.3 Deploy keys for generic Git URLs
- [ ] 5.4 Push → deploy trigger with branch filter and auto-deploy toggle
- [ ] 5.5 Commit status reporting on GitHub and GitLab

## 6. Routing and TLS

- [ ] 6.1 Full route config push to Caddy and re-push after Caddy restart
- [ ] 6.2 Wildcard certificate via Cloudflare DNS-01
- [ ] 6.3 Custom domains with DNS verification before issuance
- [ ] 6.4 HTTPS redirect, www preference and custom redirects
- [ ] 6.5 Basic auth and IP allowlist per service
- [ ] 6.6 Panel subdomain and self-signed setup port

## 7. Config, logs and metrics

- [ ] 7.1 Master key file, data key, AES-256-GCM variable encryption
- [ ] 7.2 Environment-level variables and overrides
- [ ] 7.3 `${{service.NAME}}` references with validation
- [ ] 7.4 `.env` import and pending changes with apply-restart
- [ ] 7.5 Log collector with secret masking and SSE fan-out
- [ ] 7.6 Daily zstd log files with time index, retention and search
- [ ] 7.7 Build logs tied to image retention
- [ ] 7.8 Metrics sampler (Docker stats + host), raw and hourly tables, retention
- [ ] 7.9 HTTP metrics from Caddy access logs; certificate expiry tracking

## 8. Panel

- [ ] 8.1 Vite + React + TS strict, Tailwind, shadcn/ui, embedded in the binary
- [ ] 8.2 Generated API client, TanStack Query and Router, Lingui (en, es)
- [ ] 8.3 Setup wizard (admin, domain, Git provider, sample app)
- [ ] 8.4 Login, 2FA and sessions screens
- [ ] 8.5 Home: project cards with status and server usage strip
- [ ] 8.6 Project canvas with services and their connections
- [ ] 8.7 Service pages: deploys, variables, domains, settings
- [ ] 8.8 Deploy page: steps, commits, live build log, failure explanation
- [ ] 8.9 Logs viewer with live mode and search
- [ ] 8.10 Metrics charts with uPlot
- [ ] 8.11 Command palette, shortcuts, favorites, recent activity
- [ ] 8.12 Responsive layout for phones, WCAG AA pass, performance budget in CI
- [ ] 8.13 Vitest setup and Playwright tests for login, create service, view deploy

## 9. CLI

- [ ] 9.1 CLI skeleton sharing the generated API client, `--json` everywhere
- [ ] 9.2 Device-flow login and token env var
- [ ] 9.3 `deploy` with streamed progress, `logs` with follow
- [ ] 9.4 `env` list/set/unset/import and `rollback`

## 10. Install and release

- [ ] 10.1 Install script: OS/arch checks, Docker install, containers up, master key display
- [ ] 10.2 Firewall prompt with DOCKER-USER rules
- [ ] 10.3 Socket proxy with allowlisted endpoints
- [ ] 10.4 Uninstall command with keep-apps option
- [ ] 10.5 GoReleaser, cosign signatures, SBOM, GHCR + GitLab registry
- [ ] 10.6 Installer signature verification
- [ ] 10.7 CLI install script (sh, PowerShell), Homebrew, Scoop/winget
- [ ] 10.8 Security scans in CI: govulncheck, gitleaks, gosec, CodeQL, Trivy
- [ ] 10.9 Idle memory check against the 300 MB budget
- [ ] 10.10 Sample app in `examples/`

## 11. Docs and v0.1 release

- [ ] 11.1 Starlight site with landing, install guide and concepts (en, es)
- [ ] 11.2 Page per feature (part of each feature's definition of done)
- [ ] 11.3 Full install test on a fresh Multipass VM, then on the VPS
- [ ] 11.4 60-second demo video and screenshots
- [ ] 11.5 Tag v0.1.0
