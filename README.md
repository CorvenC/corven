<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/brand/logo-dark.svg">
    <img src="docs/brand/logo-light.svg" width="72" height="72" alt="Corven logo">
  </picture>
</p>

<h1 align="center">Corven</h1>

<p align="center"><strong>Your own cloud, on your own server.</strong><br>
Self-hosted, open source deployment platform. Push to deploy on your own VPS.</p>

<p align="center"><a href="#en-español">Leer en español</a></p>

---

> **Status: planning.** No code yet. The first release (v0.1) is planned for the end of December 2026. Watch the repo to follow along.

## What it is

Corven gives a single VPS the comfort of Vercel or Railway:

- **One-command install** on Ubuntu or Debian (amd64 and arm64).
- **Push to deploy** from GitHub or GitLab, for any app with a Dockerfile or any prebuilt image.
- **Zero-downtime deploys** with health checks, and **one-click rollback** to any of the last 10 versions.
- **Automatic HTTPS** on subdomains or your own domains, even without a domain (sslip.io).
- **Encrypted environment variables**, shared per environment, with references between services.
- **Live logs** with 7 days of searchable history, and **metrics** for the server and each app.
- A **public API** and a **CLI** for Windows, macOS and Linux.

## Why another one

- **Security first:** Docker is reached through a socket proxy, secrets are encrypted at rest, releases are signed, and 2FA is built in.
- **Light:** the platform uses under 300 MB of RAM, so a 2 GB server still has room for your apps.
- **Feels like Railway:** a canvas of connected services, readable deploy pages, and errors that explain themselves.
- **English and Spanish** from day one.

## Planned install

```sh
curl -fsSL https://gitlab.com/corvenc/corven/-/raw/main/install.sh | bash
```

This does not work yet.

## Project docs

- [Roadmap](docs/ROADMAP.md)
- [Specs for v0.1](openspec/changes/add-v0-1-core/proposal.md) (OpenSpec)
- [Architecture decisions](docs/adr/README.md)
- [Brand](docs/brand/BRAND.md)

## Contributing

The repository is public from day one, but large pull requests are not accepted until v0.1 is out. Issues and ideas are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

The canonical repository is [gitlab.com/corvenc/corven](https://gitlab.com/corvenc/corven), mirrored to GitHub.

## License

[AGPL-3.0](LICENSE). You can use, modify and self-host Corven for free. If you offer a modified version as a network service, you must publish your changes.

---

## En español

**Corven** es una plataforma de deploy open source que instalás en tu propio servidor. Conectás un repo, hacés push y tu app sale con HTTPS, sin caídas durante el deploy, con logs y con rollback. Es como tener tu propio Vercel o Railway en una VPS, sin pagarlos.

> **Estado: planificación.** Todavía no hay código. La primera versión (v0.1) está planeada para fines de diciembre de 2026.

El panel y la documentación van a estar completos en español desde el primer día. Licencia AGPL-3.0.
