# 0002. Caddy as a separate container driven by its admin API

Date: 2026-11-21 · Status: Accepted

## Context
The platform needs a reverse proxy with automatic TLS that can change routes without reloads. User apps must stay online while the platform restarts or upgrades.

## Decision
Caddy runs in its own container. The platform configures it only through Caddy's admin API on an internal network and re-pushes the full config whenever Caddy restarts.

## Consequences
- Platform restarts and upgrades do not drop app traffic.
- Automatic HTTPS, wildcard certificates through DNS-01 modules, and hot route changes come for free.
- The Caddy image must include the Cloudflare DNS module, so we publish our own Caddy build.

## Alternatives considered
- Caddy embedded as a Go library: one process less, but a platform crash takes every app down.
- Traefik: configures itself from Docker labels, but is harder to debug and less direct to drive from code.
- Nginx: TLS and reloads would have to be built by hand.
