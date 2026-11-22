# Corven brand

## Name
- Product: **Corven**. CLI command and binary: `corven`.
- Tagline: *Your own cloud, on your own server.*
- Written "Corven" in prose, `corven` in code, commands and URLs.
- GitLab (primary): `gitlab.com/corvenc/corven`. GitHub (mirror): `github.com/CorvenC/corven`.
- Repository config file: `corven.toml`.
- No domain yet; docs on GitLab/GitHub Pages until one is bought.

Name checks (2026-11-21): no deployment or DevOps product uses the name; on GitHub the closest is a 1-star CSS tool. Known overlaps in other industries: Grupo Corven (motorcycles and auto parts, Argentina) and similarly named software companies (Corvena, Corvenia). `gitlab.com/corven` and `github.com/corven` belong to existing users, hence `corvenc` on GitLab and `CorvenC` on GitHub. Earlier candidates were rejected because other deploy platforms already use them (Keelson, Cleat, Berth, Mooring).

## Logo
An isometric "C" made of three stacked blocks: containers piled up into the letter, with the middle block set back to open the C.

- [`logo.svg`](logo.svg): the official mark, flat isometric faces drawn with `currentColor` at different opacities, so it is white on dark backgrounds and black on light ones. Readable down to 16 px.
- [`logo-dark.svg`](logo-dark.svg) and [`logo-light.svg`](logo-light.svg): fixed-color versions for places where `currentColor` does not work (README images, social cards).
- [`logo-concept.webp`](logo-concept.webp): the original concept artwork (logo + wordmark on black). Reference only; use the SVG in products and docs.
- Wordmark: logo + "Corven" in Geist 600.

## Color
Monochrome brand. No brand accent color.

| Token | Light | Dark | Use |
|---|---|---|---|
| `--bg` | `#ffffff` | `#0a0a0a` | page background |
| `--surface` | `#f5f5f5` | `#141414` | cards, panels |
| `--line` | `#e5e5e5` | `#262626` | borders |
| `--muted` | `#737373` | `#a3a3a3` | secondary text |
| `--fg` | `#0a0a0a` | `#fafafa` | text, logo, primary buttons |

Status colors are only for state, never for branding: success `#16a34a`, warning `#d97706`, error `#dc2626` (lighter variants in dark mode).

## Type
- Interface and headings: **Geist** (400, 500, 600, 700).
- Logs, code, IDs and numbers in tables: **JetBrains Mono** (400, 500).
- Both from Google Fonts, with system fallbacks.
