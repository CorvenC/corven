# 0005. Envelope encryption for secrets

Date: 2026-11-21 · Status: Accepted

## Context
Apps store API keys and passwords in environment variables. A leaked database dump or backup must not reveal them.

## Decision
Values are encrypted with AES-256-GCM using a data key. The data key is stored in the database encrypted by a master key. The master key lives in a file on the host, outside the database and outside backups, and is shown once at install so the operator can keep a copy.

## Consequences
- A database dump or backup alone reveals nothing.
- Rotating the master key only re-encrypts the data key.
- Losing the master key makes the secrets unrecoverable on another server: the installer and the docs warn about it clearly.

## Alternatives considered
- Plain text in the database: unacceptable for this product.
- Key derived from the admin password: a password change or a forgotten password breaks everything.
- External secret managers (Vault): too heavy for the target user.
