# Roadmap

Product: **Corven** (see [brand/BRAND.md](brand/BRAND.md)). Repo: `gitlab.com/corvenc/corven`, mirrored to `github.com/CorvenC/corven`.

## v0.1 — target 2026-12-31
Change: `openspec/changes/add-v0-1-core`. One-command install, admin with 2FA, GitHub/GitLab integration, Dockerfile and image deploys with zero downtime, rollback, subdomains and TLS (Cloudflare wildcard or sslip.io), custom domains, encrypted variables with references, live and stored logs, basic metrics, public API with scoped tokens, basic CLI.

## v1.0 — no fixed date
Each item becomes its own OpenSpec change.

- Preview environments per merge request: `pr-12--app.domain`, fresh Postgres per preview, protected by default, auto-removed on merge or after 7 idle days, limit per project, comment on the MR.
- Managed Postgres and Redis: official versions, persistent volumes, connection URL injected, tunnel through the CLI, major-version upgrade assistant.
- Backups: platform database and app databases, daily (7 daily + 4 weekly), local or S3-compatible, encrypted, one-click restore, weekly automatic restore test.
- Workers and cron jobs (same image, different command).
- docker-compose import into platform services.
- Around 10 app templates (WordPress, n8n, Uptime Kuma, Plausible...).
- Notifications: Discord, Telegram, email; per-project rules; in-panel notification center; signed outgoing webhooks.
- Audit log of every important action.
- Platform self-update from the panel with backup and automatic rollback.
- Export/import the whole platform to move to another server.
- Optional AI features (off by default): failure diagnosis and Dockerfile generation with Ollama, Claude or OpenAI.
- Community: Discord, accept external contributions, GitHub PR bridge.
- Public announcement (Reddit, Hacker News, X) when v1.0 ships. Donations decided then.

## After v1.0
Multi-server (agent per node, own registry), teams and roles, web terminal into containers, rate limiting per app, log drains, SQL console, passkeys, larger template catalog.

## Never
Kubernetes, serverless functions, Windows servers, telemetry.
