# Contributing to Corven

Thanks for your interest. Corven is in its planning stage, so here is what helps right now and how the project works.

## Right now (before v0.1)

- **Issues are welcome:** bugs in the specs, ideas, questions, use cases.
- **Large pull requests are not accepted yet.** The architecture is still moving fast. Small fixes to docs and typos are fine.
- External code contributions open with v0.1.

## Where things live

- The canonical repository is on GitLab: [gitlab.com/corvenc/corven](https://gitlab.com/corvenc/corven). GitHub is a mirror.
- Pull requests opened on GitHub are brought into GitLab as branches, keeping your authorship.
- Design and behavior are written first as [OpenSpec](https://github.com/Fission-AI/OpenSpec) changes in `openspec/`. Important decisions are recorded in `docs/adr/`.

## How changes are made

- **Trunk-based:** `main` always works. Each change is a short branch with a merge request.
- **Small merge requests:** one thing per MR, about 400 changed lines at most.
- **Commit messages** follow [Conventional Commits](https://www.conventionalcommits.org/) with a scope, for example `feat(deploy): retry health checks`.
- **Sign your commits off** with the [Developer Certificate of Origin](https://developercertificate.org/): use `git commit -s`, which adds a `Signed-off-by` line.
- **Code, comments and commit messages are in English.** User-facing docs exist in English and Spanish.

## Definition of done

A change is merged when:

1. CI is green (lint, tests, build, security scans).
2. It is documented (English and Spanish user docs where it applies).
3. It was tested on a clean Ubuntu install, not only in development.
4. A maintainer reviewed and understood the diff.

## Development setup

Instructions will be added with the first code. The expected setup is Linux (or WSL2 on Windows) with Docker Engine, Go, Node with pnpm, and [Task](https://taskfile.dev).

## Code of conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md).
