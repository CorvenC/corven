# 0007. AGPL-3.0 license

Date: 2026-11-21 · Status: Accepted

## Context
The project is open source and self-hosted. Competitors are also open source, so openness alone is not a differentiator. The risk is a company taking the code and selling it as a closed hosted service.

## Decision
The project is licensed under AGPL-3.0. Contributions require a DCO sign-off (`Signed-off-by`), not a CLA.

## Consequences
- Anyone can use, modify and self-host it for free.
- Anyone who offers a modified version as a network service must publish their changes.
- Some companies avoid AGPL dependencies; acceptable, since the product is an application, not a library.

## Alternatives considered
- MIT or Apache-2.0: maximum adoption, but allows closed hosted forks.
- CLA: would allow relicensing later, but drives contributors away.
