# 0003. Huma on chi, OpenAPI generated from code

Date: 2026-11-21 · Status: Accepted

## Context
The API is public and is the only interface used by the panel and the CLI. Its documentation, validation and TypeScript client must never drift apart.

## Decision
Endpoints are defined with Huma on top of the chi router. Go structs are the single source of truth: Huma generates the OpenAPI document, validates requests and returns RFC 9457 errors. The panel's TypeScript client is generated from that document in CI.

## Consequences
- Changing an endpoint updates the docs, the validation and, after regeneration, the client.
- A broken contract between back and front fails at compile time in the panel.

## Alternatives considered
- Spec-first with oapi-codegen: more control, a lot of hand-written YAML.
- chi with hand-written validation: the OpenAPI document would fall out of date.
