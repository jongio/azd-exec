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
```

## Deployment

The site is automatically deployed to GitHub Pages on push to main via GitHub Actions.

**Live Site**: https://jongio.github.io/azd-exec/

## Documentation

- [CLI Reference](../cli/docs/cli-reference.md) - Complete command reference
- [Security Review](../cli/docs/SECURITY-REVIEW.md) - Security documentation
- [Threat Model](../cli/docs/THREAT-MODEL.md) - Security threat analysis

## Structure

- `src/components/` - Reusable Astro components
- `src/pages/` - Site pages (auto-routed)
- `src/styles/` - Global CSS and Tailwind imports
- `public/` - Static assets
- `scripts/` - Build-time generation scripts
