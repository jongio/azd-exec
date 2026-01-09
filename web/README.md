---
title: azd-exec Documentation Website
description: Astro-based documentation website for azd-exec extension
lastUpdated: 2026-01-09
tags: [website, documentation, astro]
---

# azd-exec Documentation Website

Documentation website for the azd-exec extension, built with Astro and Tailwind CSS.

## Development

```bash
# Install dependencies
pnpm install

# Start dev server
pnpm run dev

# Build for production
pnpm run build

# Preview production build
pnpm run preview

# Generate changelog page
pnpm run generate:changelog

# Generate OG image for social media
pnpm run generate:og-image
```

## Deployment

The site is automatically deployed to GitHub Pages on push to main via GitHub Actions.

**Live Site**: https://jongio.github.io/azd-exec/

## Documentation

- [CLI Reference](../cli/docs/cli-reference.md) - Complete command reference
- [Security Review](../cli/docs/security-review.md) - Security documentation
- [Threat Model](../cli/docs/threat-model.md) - Security threat analysis

## Structure

- `src/components/` - Reusable Astro components
- `src/pages/` - Site pages (auto-routed)
- `src/styles/` - Global CSS and Tailwind imports
- `public/` - Static assets
- `scripts/` - Build-time generation scripts
