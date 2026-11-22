# 0004. Docker access through a socket proxy

Date: 2026-11-21 · Status: Accepted

## Context
Controlling Docker requires its socket, and access to the socket is equivalent to root on the host. Mounting it directly into a web-facing container means a panel vulnerability is a full server compromise. Most competing tools do exactly that.

## Decision
The platform reaches Docker only through a socket proxy that allows the specific API endpoints it needs and denies the rest (for example, no privileged containers and no host bind mounts outside the platform's own data paths).

## Consequences
- A compromised platform process cannot freely take over the host.
- Every new Docker call the platform needs must be added to the allowlist on purpose.
- Security becomes a real differentiator we can explain in the docs.

## Alternatives considered
- Mount the socket directly: simplest, root-equivalent exposure.
- Rootless Docker: better isolation, but breaks common setups and complicates the installer for v0.1.
