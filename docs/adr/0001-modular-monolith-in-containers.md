# 0001. Modular monolith, everything in containers

Date: 2026-11-21 · Status: Accepted

## Context
One developer, a small server (2 GB RAM), and a product whose main promise is "install with one command". Multi-server support is wanted later, not now.

## Decision
The platform is a single Go binary with strictly separated internal modules (auth, git, build, deploy, proxy, logs, metrics, api). It runs as a container next to its own PostgreSQL and Caddy containers.

## Consequences
- One artifact to build, sign, ship and upgrade; low memory use.
- Install, upgrade and uninstall are clean because nothing is installed on the host besides Docker.
- Module boundaries must be respected so a per-server agent can be extracted for multi-server.

## Alternatives considered
- Control plane + agent per server from day one: doubles the work before there is a second server.
- Microservices: operational cost far beyond one developer.
- Binary under systemd: lighter, but scatters files across the host.
