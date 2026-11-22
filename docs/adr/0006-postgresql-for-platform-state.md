# 0006. PostgreSQL for platform state and jobs

Date: 2026-11-21 · Status: Accepted

## Context
The platform stores projects, deploys, sessions, metrics and a durable job queue. SQLite would be lighter for a single server.

## Decision
Platform state lives in a dedicated PostgreSQL container, never shared with user apps. The job queue (River) and metrics tables live in the same database.

## Consequences
- One database for state, jobs and metrics; jobs survive restarts.
- Ready for multi-server, where several nodes need a shared state.
- Costs around 100 MB of RAM; Postgres is tuned for small servers and checked against the 300 MB budget.

## Alternatives considered
- SQLite: lighter and simpler backups, but a harder path to multi-server and to a durable queue.
- Shared Postgres with user apps: saves RAM, but an app can break or read the platform's data.
