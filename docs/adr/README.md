# Architecture Decision Records

Short records of important decisions and why they were made. One file per decision, never edited after acceptance: a change of mind gets a new ADR that supersedes the old one.

| # | Decision | Status |
|---|---|---|
| [0001](0001-modular-monolith-in-containers.md) | Modular monolith, everything in containers | Accepted |
| [0002](0002-caddy-as-separate-container.md) | Caddy as a separate container driven by its admin API | Accepted |
| [0003](0003-huma-for-the-api.md) | Huma on chi, OpenAPI generated from code | Accepted |
| [0004](0004-docker-socket-proxy.md) | Docker access through a socket proxy | Accepted |
| [0005](0005-secrets-encryption.md) | Envelope encryption for secrets | Accepted |
| [0006](0006-postgresql-for-platform-state.md) | PostgreSQL for platform state and jobs | Accepted |
| [0007](0007-agpl-license.md) | AGPL-3.0 license | Accepted |
| [0008](0008-gitlab-primary-github-mirror.md) | GitLab as primary, GitHub as mirror | Accepted |

Template: Context → Decision → Consequences → Alternatives considered.
