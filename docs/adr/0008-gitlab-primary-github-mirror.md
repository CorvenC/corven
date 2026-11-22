# 0008. GitLab as primary, GitHub as mirror

Date: 2026-11-21 · Status: Accepted

## Context
The maintainer works in GitLab, but the open source community and contribution graphs live on GitHub.

## Decision
The canonical repository, CI, issues board and milestones live on GitLab. Every push goes to both: the local `origin` remote has two push URLs, GitLab and GitHub, so the mirror never drifts and no access token has to be stored in GitLab. A GitLab push mirror can replace this later if needed. Pull requests opened on GitHub are brought into GitLab as branches (keeping authorship) and merged there; issues opened on GitHub are synced to the GitLab board.

## Consequences
- The maintainer keeps one workflow; the project is discoverable on GitHub and commits count there.
- Bringing GitHub PRs into GitLab needs a small script or bot.
- Community chat happens on Discord.

## Alternatives considered
- GitHub as primary: simpler for contributors, changes the maintainer's workflow.
- GitLab only: invisible to most of the community.
