# azd-exec Documentation Website

## Overview

Create a documentation website for azd-exec extension using Astro, styled with Tailwind CSS, hosted on GitHub Pages. The site provides comprehensive documentation, examples, and reference material for the extension.

## Goals

1. Provide clear, searchable documentation for azd-exec
2. Offer copy-paste ready installation and usage examples
3. Generate CLI reference from source
4. Auto-generate changelog pages from cli/CHANGELOG.md
5. Support dark mode and mobile-responsive design
6. Deploy automatically via GitHub Actions

## Scope

### In Scope

- **Home Page**: Overview, features, quick start, security notice
- **Getting Started**: Detailed installation and first-time setup
- **CLI Reference**: Command documentation (run, version)
- **Examples**: Common use cases with full scripts
- **Changelog**: Full version history
- **Components**: Reusable UI elements (code blocks, install tabs, layout)
- **Auto-generation**: CLI reference and changelog from source
- **GitHub Pages**: Static deployment with automated builds

### Out of Scope

- Interactive dashboard (azd-app specific)
- Guided tour (azd-app specific)
- MCP integration documentation (azd-exec doesn't use MCP)
- E2E tests (defer to future iteration)
- Search functionality (defer to future iteration)
- Multi-language support (English only)

## Architecture

### Technology Stack

- **Framework**: Astro 4.x (static site generation)
- **Styling**: Tailwind CSS 4.x
- **Code Highlighting**: astro-expressive-code
- **Build**: TypeScript, tsx for scripts
- **Deployment**: GitHub Pages (static)

### Directory Structure

```
web/
├── public/                     # Static assets
│   ├── favicon.svg
│   └── screenshots/            # Future: command execution screenshots
├── src/
│   ├── components/
│   │   ├── Layout.astro        # Main layout wrapper
│   │   ├── Header.astro        # Top navigation
│   │   ├── Footer.astro        # Bottom navigation
│   │   ├── CodeBlock.astro     # Code with copy button
│   │   └── InstallTabs.astro   # Platform-specific install
│   ├── pages/
│   │   ├── index.astro         # Home page
│   │   ├── getting-started.astro
│   │   ├── examples.astro
│   │   ├── reference/
│   │   │   ├── cli/
│   │   │   │   ├── index.astro       # CLI overview
│   │   │   │   ├── run.astro         # Generated
│   │   │   │   └── version.astro     # Generated
│   │   │   └── changelog/
│   │   │       └── index.astro       # Generated
│   │   └── security.astro
│   ├── styles/
│   │   └── global.css          # Global styles, Tailwind imports
│   └── env.d.ts                # TypeScript types
├── scripts/
│   ├── generate-cli-reference.ts    # Parse CLI docs to pages
│   └── generate-changelog.ts        # Parse CHANGELOG.md to pages
├── astro.config.mjs
├── tailwind.config.ts
├── tsconfig.json
├── package.json
└── README.md
```

## Pages

### Home Page (`index.astro`)

**Purpose**: Landing page with overview and quick start

**Sections**:
- Hero: "Run Scripts with Azure Context" + key value prop
- Security Warning: Prominent callout about credential access
- Features Grid: 6 key features with icons
- Quick Start: 3-step installation process
- Example Script: Bash and PowerShell examples
- Footer: Links to GitHub, docs, security

**Design**:
- Clean, modern, professional
- Security notice in yellow/orange callout
- Code blocks with platform tabs (Windows/macOS/Linux)
- Responsive grid layout

### Getting Started (`getting-started.astro`)

**Purpose**: Detailed first-time setup guide

**Sections**:
1. Prerequisites (azd installed)
2. Enable Extensions & Install
3. Verify Installation
4. Run Your First Script
5. Script Security Best Practices
6. Next Steps (links to examples, CLI reference)

### CLI Reference (`reference/cli/*.astro`)

**Purpose**: Complete command documentation

**Commands to Document**:
- `azd exec` - Execute scripts with azd context
- `azd exec version` - Show version

**Auto-generation**:
- Read from `cli/docs/commands/*.md` (future)
- For now, manually create pages
- Include: syntax, flags, examples, common errors

### Examples (`examples.astro`)

**Purpose**: Real-world usage patterns

**Examples**:
1. **Azure Deployment Script**: Deploy resources with azd context
2. **Environment Variable Access**: Read azd env vars in script
3. **Multi-Step Build**: Chained commands with arguments
4. **Cross-Platform Script**: Works on Windows and Linux
5. **Interactive Setup**: Script with user prompts
6. **Error Handling**: Proper exit codes and logging

**Format**: Each example shows:
- Use case description
- Complete script code
- How to run it
- Expected output
- Security considerations

### Changelog (`reference/changelog/index.astro`)

**Purpose**: Full version history

**Auto-generation**:
- Parse `cli/CHANGELOG.md`
- Group by version
- Show date, changes, commit links
- Support conventional commit types (feat, fix, docs, chore)

**Design**:
- Reverse chronological
- Version badges
- Clickable commit hashes
- Release notes links

### Security (`security.astro`)

**Purpose**: Security best practices and threat model

**Sections**:
- What Scripts Can Access
- Safe Practices
- Dangerous Practices
- Threat Model Overview
- Link to full docs (security-review.md, threat-model.md)

## Components

### Layout.astro

**Purpose**: Main wrapper for all pages

**Features**:
- HTML head with meta tags, title, description
- Header component
- Main content area
- Footer component
- Dark mode toggle (future)
- Breadcrumb navigation (optional)

### Header.astro

**Purpose**: Top navigation bar

**Elements**:
- Logo/title
- Navigation links: Home, Getting Started, Examples, CLI Reference, Changelog
- GitHub icon link
- Responsive hamburger menu (mobile)

### Footer.astro

**Purpose**: Bottom navigation and links

**Elements**:
- Copyright notice
- Quick links: Security, Contributing, License
- Social links: GitHub
- Built with Astro badge

### CodeBlock.astro

**Purpose**: Syntax-highlighted code with copy button

**Features**:
- Language detection
- Line numbers (optional)
- Copy to clipboard button
- Theme integration (light/dark)
- Inline vs block modes

### InstallTabs.astro

**Purpose**: Platform-specific installation commands

**Features**:
- Tabs for Windows, macOS, Linux
- Code blocks for each platform
- Default to user's OS (client-side JS)
- Copy buttons

## Scripts

### generate-cli-reference.ts

**Purpose**: Generate CLI reference pages from docs

**Process**:
1. Read `cli/docs/commands/*.md`
2. Parse markdown for command syntax, flags, examples
3. Generate `src/pages/reference/cli/[command].astro`
4. Create index page with all commands

**Future**: Once CLI docs are structured

### generate-changelog.ts

**Purpose**: Generate changelog page from CHANGELOG.md

**Process**:
1. Read `cli/CHANGELOG.md`
2. Parse version sections, dates, changes
3. Extract commit hashes and PR numbers
4. Generate `src/pages/reference/changelog/index.astro`
5. Create links to GitHub commits/PRs

**Format**: Similar to azd-app implementation

## Styling

### Theme

**Colors**:
- Primary: Azure blue (#0078D4)
- Success: Green (#107C10)
- Warning: Orange/Yellow (#F59E0B)
- Error: Red (#DC2626)
- Neutral: Grays for text and backgrounds

**Typography**:
- Font: System font stack (Segoe UI, Roboto, sans-serif)
- Headings: Bold, larger sizes
- Code: Monospace (Consolas, Monaco, Courier New)

**Dark Mode**:
- Support light and dark themes
- Use CSS variables or Tailwind dark: prefix
- Toggle in header (future)

### Responsive Design

- **Mobile**: Single column, hamburger menu
- **Tablet**: Adjusted spacing, compact nav
- **Desktop**: Full layout, sidebar (future)

## Build & Deployment

### Build Process

1. Install dependencies: `pnpm install`
2. Generate pages: Run scripts (changelog, CLI reference)
3. Build site: `pnpm run build`
4. Output: `dist/` directory

### NPM Scripts

```json
{
  "scripts": {
    "dev": "astro dev",
    "build": "tsx scripts/generate-changelog.ts && astro build",
    "preview": "astro preview",
    "generate:changelog": "tsx scripts/generate-changelog.ts"
  }
}
```

### GitHub Actions

**Workflow**: `.github/workflows/docs.yml`

**Trigger**: Push to main, manual dispatch

**Steps**:
1. Checkout code
2. Setup Node.js (20.x)
3. Install pnpm
4. Install dependencies (`pnpm install --dir web`)
5. Generate pages
6. Build site (`pnpm --dir web run build`)
7. Deploy to GitHub Pages

**Pages Setup**:
- Source: GitHub Actions
- Branch: Not applicable (Actions deployment)
- Base URL: `/azd-exec/`

### Local Development

```bash
cd web
pnpm install
pnpm run dev
# Visit http://localhost:4321/azd-exec/
```

## Configuration

### astro.config.mjs

```js
import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';
import expressiveCode from 'astro-expressive-code';

export default defineConfig({
  site: 'https://jongio.github.io/azd-exec/',
  base: '/azd-exec/',
  integrations: [expressiveCode()],
  vite: {
    plugins: [tailwindcss()]
  },
  output: 'static'
});
```

### package.json

**Dependencies**:
- astro: ^4.x
- @astrojs/mdx: ^3.x
- astro-expressive-code: ^0.x
- tailwindcss: ^4.x
- @tailwindcss/vite: ^4.x
- typescript: ^5.x
- tsx: ^4.x

**DevDependencies**: (none specific)

## Content Strategy

### Tone & Voice

- **Professional**: Clear, technical, precise
- **Helpful**: Anticipate user questions
- **Security-conscious**: Emphasize safety
- **Concise**: Get to the point quickly

### Security Messaging

**Prominent Warnings**:
- Home page: Large callout about credential access
- Getting Started: Repeat security practices
- Examples: Include security notes for each script
- Dedicated Security page

**Best Practices**:
- Always review scripts before running
- Use HTTPS for downloads
- Verify sources
- Avoid piping unknown scripts
- Store secrets in Key Vault, not env vars

## Future Enhancements

1. **Search**: Algolia or static search
2. **Versioned Docs**: Support multiple azd-exec versions
3. **Dark Mode Toggle**: User preference + system detection
4. **Screenshots**: Terminal output examples
5. **Interactive Examples**: Copy and try in browser
6. **E2E Tests**: Playwright for visual regression
7. **Analytics**: Google Analytics or privacy-focused alternative
8. **Feedback Widget**: Allow users to rate docs
9. **CLI Docs Structure**: Formalize markdown format for auto-generation

## Success Metrics

1. **Deployment**: Site builds and deploys successfully
2. **Performance**: Lighthouse score >90
3. **Accessibility**: WCAG AA compliance
4. **Mobile**: Responsive on all screen sizes
5. **SEO**: Proper meta tags, sitemap
6. **Usage**: Track page views (future)

## Implementation Timeline

**Phase 1**: Core Structure (Week 1)
- Setup web/ directory
- Create Layout, Header, Footer
- Build home page
- Configure deployment

**Phase 2**: Content (Week 2)
- Getting Started page
- CLI Reference pages
- Examples page
- Security page

**Phase 3**: Automation (Week 3)
- generate-changelog.ts script
- Build integration
- GitHub Actions workflow

**Phase 4**: Polish (Week 4)
- Dark mode
- Responsive design refinements
- Performance optimization
- SEO improvements

## References

- azd-app web structure: `https://github.com/jongio/azd-app/tree/main/web`
- Astro docs: `https://docs.astro.build`
- Tailwind CSS: `https://tailwindcss.com`
- GitHub Pages: `https://docs.github.com/en/pages`
