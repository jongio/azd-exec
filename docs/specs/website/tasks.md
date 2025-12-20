# azd-exec Website Tasks

## TODO: Setup Project Foundation

**Description**: Initialize web/ directory with Astro, package manager, and base configuration files.

**Tasks**:
1. Create `web/` directory structure:
   ```
   web/
   ├── public/
   ├── src/
   │   ├── components/
   │   ├── pages/
   │   └── styles/
   ├── scripts/
   ├── package.json
   ├── astro.config.mjs
   ├── tailwind.config.ts
   ├── tsconfig.json
   └── README.md
   ```

2. Create `web/package.json`:
   - Name: `@jongio/azd-exec-docs`
   - Scripts: dev, build, preview, generate:changelog
   - Dependencies: astro, @astrojs/mdx, astro-expressive-code, tailwindcss, @tailwindcss/vite
   - DevDependencies: typescript, tsx, @types/node

3. Create `web/astro.config.mjs`:
   - Site: `https://jongio.github.io/azd-exec/`
   - Base: `/azd-exec/`
   - Integrations: expressiveCode
   - Vite: tailwindcss plugin
   - Output: static

4. Create `web/tsconfig.json`:
   - Extends astro/tsconfigs/strict
   - CompilerOptions for types
   - Include src/**

5. Create `web/tailwind.config.ts`:
   - Content: ./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}
   - Theme extensions for custom colors
   - Plugins: none initially

6. Create `web/src/env.d.ts`:
   - Reference .astro/types.d.ts

7. Create `web/src/styles/global.css`:
   - Tailwind imports: @tailwind base, components, utilities
   - Custom CSS variables for colors
   - Base styles

8. Create `web/public/favicon.svg`:
   - Simple SVG icon or copy from azd-app

9. Create `web/README.md`:
   - Purpose: Documentation website for azd-exec
   - Development commands
   - Build and deployment info

**Acceptance**:
- `pnpm install --dir web` succeeds
- `pnpm --dir web run dev` starts dev server
- Browser shows blank Astro page at http://localhost:4321/azd-exec/

---

## TODO: Build Layout Component

**Description**: Create main Layout.astro wrapper used by all pages.

**Tasks**:
1. Create `web/src/components/Layout.astro`:
   - Props: title, description, contentClass (optional)
   - HTML head with:
     - Meta charset, viewport
     - Title format: "{title} | azd exec"
     - Description meta tag
     - Favicon link
   - Body structure:
     - Header component
     - Main content slot
     - Footer component
   - Global CSS import
   - Base URL handling for GitHub Pages

2. Add TypeScript interface for props

3. Test with minimal page

**Acceptance**:
- Layout renders with proper meta tags
- Base URL correctly applied to links
- Title and description configurable

---

## TODO: Build Header Component

**Description**: Top navigation bar with links and responsive design.

**Tasks**:
1. Create `web/src/components/Header.astro`:
   - Logo/title: "azd exec"
   - Navigation links:
     - Home (/)
     - Getting Started (/getting-started)
     - Examples (/examples)
     - CLI Reference (/reference/cli/)
     - Changelog (/reference/changelog/)
   - GitHub icon link (external)
   - Responsive hamburger menu (mobile)

2. Style with Tailwind:
   - Fixed or sticky header
   - Horizontal nav on desktop
   - Hamburger menu on mobile
   - Hover states
   - Active link highlighting

3. Add client-side JS for mobile menu toggle

**Acceptance**:
- All navigation links work
- Responsive on mobile/tablet/desktop
- GitHub link opens in new tab
- Active page highlighted

---

## TODO: Build Footer Component

**Description**: Bottom section with links and copyright.

**Tasks**:
1. Create `web/src/components/Footer.astro`:
   - Copyright notice: "© 2025 azd exec"
   - Quick links:
     - Security (/security)
     - Contributing (GitHub CONTRIBUTING.md)
     - License (GitHub LICENSE)
   - Social/External links:
     - GitHub repo
     - azd documentation
   - Built with Astro badge

2. Style with Tailwind:
   - Dark background
   - Light text
   - Grid layout for links
   - Responsive columns

**Acceptance**:
- Footer renders on all pages
- Links functional
- Responsive layout

---

## TODO: Build CodeBlock Component

**Description**: Reusable code block with syntax highlighting and copy button.

**Tasks**:
1. Create `web/src/components/CodeBlock.astro`:
   - Props: code, language, showLanguage, showCopy, showLineNumbers
   - Use astro-expressive-code's Code component
   - Add copy button overlay
   - Optional language badge

2. Configure expressive-code:
   - Create `web/ec.config.mjs`
   - Themes: github-light, github-dark
   - Theme selector: [data-theme="${theme.type}"]
   - Style overrides: border radius, padding

3. Add copy-to-clipboard functionality (client JS)

**Acceptance**:
- Syntax highlighting works
- Copy button functional
- Language indicator shows
- Dark mode compatible (future)

---

## TODO: Build InstallTabs Component

**Description**: Platform-specific installation commands with tabs.

**Tasks**:
1. Create `web/src/components/InstallTabs.astro`:
   - Tabs: Windows, macOS, Linux
   - Each tab shows install commands:
     - Windows: winget or PowerShell
     - macOS: brew
     - Linux: curl script
   - Default to user's OS (client-side JS)

2. Tab switching logic:
   - Active tab highlighted
   - Content switches on click
   - Keyboard navigation (a11y)

3. Use CodeBlock for each command

**Acceptance**:
- Tabs switch correctly
- Install commands accurate
- Copy buttons work
- Default tab selects user OS

---

## TODO: Create Home Page

**Description**: Landing page with overview and quick start.

**Tasks**:
1. Create `web/src/pages/index.astro`:
   - Import Layout
   - Sections:
     - Hero: Title, tagline, CTA button
     - Security Warning: Large callout about credentials
     - Features Grid: 6 features with icons
     - Quick Start: 3 steps with InstallTabs
     - Example Scripts: Bash and PowerShell tabs
     - Footer CTA: Link to getting started

2. Use Tailwind for styling:
   - Hero: Large centered text, gradient background
   - Features: 2-column grid on desktop
   - Security: Yellow/orange callout box
   - Code examples: CodeBlock component

3. Add SEO meta tags via Layout

**Acceptance**:
- Page loads at `/azd-exec/`
- All sections render correctly
- Security warning prominent
- Install commands copyable
- Responsive design

---

## TODO: Create Getting Started Page

**Description**: Step-by-step setup guide for first-time users.

**Tasks**:
1. Create `web/src/pages/getting-started.astro`:
   - Import Layout, CodeBlock, InstallTabs
   - Sections:
     1. Prerequisites (azd installed)
     2. Enable Extensions & Install (InstallTabs)
     3. Verify Installation (azd exec version)
     4. Run Your First Script (example script)
     5. Security Best Practices (list)
     6. Next Steps (links to examples, CLI ref)

2. Include example script:
   - Simple deploy.sh
   - Shows azd env var access
   - Security review reminder

3. Add navigation: Previous (Home), Next (Examples)

**Acceptance**:
- Step-by-step flow clear
- Code examples work
- Links to other pages functional
- Security practices emphasized

---

## TODO: Create Examples Page

**Description**: Real-world usage examples with full scripts.

**Tasks**:
1. Create `web/src/pages/examples.astro`:
   - Import Layout, CodeBlock
   - Examples (6 total):
     1. Azure Deployment Script
     2. Environment Variable Access
     3. Multi-Step Build
     4. Cross-Platform Script
     5. Interactive Setup
     6. Error Handling

2. Each example:
   - Use case description
   - Full script code (CodeBlock)
   - How to run (command)
   - Expected output (CodeBlock)
   - Security note

3. Table of contents at top

**Acceptance**:
- All 6 examples complete
- Scripts copyable
- Security notes present
- Responsive layout

---

## TODO: Create CLI Reference Index

**Description**: Overview page for CLI commands.

**Tasks**:
1. Create `web/src/pages/reference/cli/index.astro`:
   - Import Layout
   - Page title: "CLI Reference"
   - Introduction paragraph
   - Command list:
     - azd exec run (link to run.astro)
     - azd exec version (link to version.astro)
     - azd exec listen (brief note, internal)

2. Each command in list:
   - Name, description
   - Link to detail page

3. Breadcrumb: Home > CLI Reference

**Acceptance**:
- Page loads at `/azd-exec/reference/cli/`
- Links to command pages work
- Breadcrumb functional

---

## TODO: Create CLI Reference - run

**Description**: Detailed documentation for `azd exec run`.

**Tasks**:
1. Create `web/src/pages/reference/cli/run.astro`:
   - Import Layout, CodeBlock
   - Sections:
     - Syntax
     - Description
     - Flags table:
       - --shell: Shell to use
       - --cwd: Working directory
       - --interactive: Interactive mode
       - --: Script arguments separator
     - Examples (5+):
       - Basic run
       - Specify shell
       - Pass arguments
       - Set working directory
       - Interactive script
     - Common Errors
     - See Also (links to other commands)

2. Flags table:
   - Name, Type, Default, Description columns
   - Responsive table

**Acceptance**:
- Complete command documentation
- Examples copyable
- Flags table clear
- Breadcrumb: Home > CLI Reference > run

---

## TODO: Create CLI Reference - version

**Description**: Documentation for `azd exec version`.

**Tasks**:
1. Create `web/src/pages/reference/cli/version.astro`:
   - Import Layout, CodeBlock
   - Sections:
     - Syntax: `azd exec version`
     - Description: Display version info
     - Example output (CodeBlock)
     - See Also

2. Keep it brief (simple command)

**Acceptance**:
- Page complete
- Example shows actual version format
- Breadcrumb: Home > CLI Reference > version

---

## TODO: Create Security Page

**Description**: Security best practices and threat model overview.

**Tasks**:
1. Create `web/src/pages/security.astro`:
   - Import Layout
   - Sections:
     - What Scripts Can Access (credentials, env vars, filesystem)
     - Safe Practices (bulleted list)
     - Dangerous Practices (what NOT to do)
     - Threat Model Overview
     - Links to full docs:
       - cli/docs/SECURITY-REVIEW.md (GitHub)
       - cli/docs/THREAT-MODEL.md (GitHub)

2. Use callout boxes:
   - Green for safe practices
   - Red for dangerous practices

3. Emphasize script review before execution

**Acceptance**:
- Security info clear
- Links to GitHub docs work
- Callouts visually distinct

---

## TODO: Create generate-changelog Script

**Description**: Auto-generate changelog page from cli/CHANGELOG.md.

**Tasks**:
1. Create `web/scripts/generate-changelog.ts`:
   - Read `cli/CHANGELOG.md`
   - Parse markdown:
     - Extract version sections (## [x.y.z])
     - Extract dates
     - Extract changes (bullet points)
     - Extract commit hashes
   - Generate `web/src/pages/reference/changelog/index.astro`:
     - Import Layout
     - Version list (reverse chronological)
     - Each version: heading, date, changes
     - Link commit hashes to GitHub
   - Create links to release tags

2. Handle conventional commits:
   - feat: ✨
   - fix: 🐛
   - docs: 📝
   - chore: 🔧

3. Error handling for missing CHANGELOG

**Acceptance**:
- Script runs: `pnpm --dir web run generate:changelog`
- Generates `web/src/pages/reference/changelog/index.astro`
- Changelog page renders correctly
- Links to commits work

---

## TODO: Integrate Build Scripts

**Description**: Add script generation to build process.

**Tasks**:
1. Update `web/package.json` scripts:
   - `build`: `tsx scripts/generate-changelog.ts && astro build`
   - `generate:changelog`: `tsx scripts/generate-changelog.ts`

2. Test build:
   - `pnpm --dir web run build`
   - Verify `dist/` created
   - Check all pages present

3. Create `.gitignore` in web/:
   - `node_modules/`
   - `dist/`
   - `.astro/`

**Acceptance**:
- Build generates all pages
- Changelog auto-generated
- dist/ contains static site

---

## TODO: Configure GitHub Actions Deployment

**Description**: Set up automated deployment to GitHub Pages.

**Tasks**:
1. Create `.github/workflows/docs.yml`:
   - Name: Deploy Docs
   - Triggers:
     - Push to main (paths: web/**)
     - Manual workflow_dispatch
   - Jobs:
     - build-and-deploy:
       - Checkout code
       - Setup Node 20
       - Install pnpm
       - Cache pnpm store
       - Install deps: `pnpm install --dir web`
       - Build: `pnpm --dir web run build`
       - Upload artifact
       - Deploy to GitHub Pages

2. Configure GitHub Pages:
   - Settings > Pages
   - Source: GitHub Actions
   - Branch: Not needed (Actions)

3. Set environment in workflow:
   - GITHUB_TOKEN (auto-provided)

**Acceptance**:
- Workflow runs on push
- Pages deploy successfully
- Site accessible at https://jongio.github.io/azd-exec/

---

## TODO: Add Spell Checking for Web

**Description**: Extend cspell to cover web/ directory.

**Tasks**:
1. Update `cspell.json` in repo root:
   - Add `web/**/*.{astro,ts,md}` to files
   - Add web-specific words to dictionary:
     - astro, mdx, expressive-code
     - tailwindcss, pnpm
     - Any project-specific terms

2. Test: `cspell "web/**/*.{astro,ts,md}" --config cspell.json`

3. Update CI workflow to include web spell check

**Acceptance**:
- Spell check runs on web files
- No false positives
- CI includes web spell check

---

## TODO: Update Main README

**Description**: Add link to documentation website.

**Tasks**:
1. Edit `README.md`:
   - Add "Documentation" section near top
   - Link: `📚 [Full Documentation](https://jongio.github.io/azd-exec/)`
   - Update Quick Start to reference website for detailed guides

2. Keep README concise, point to website for details

**Acceptance**:
- Link to docs present
- README not duplicating website content

---

## TODO: Add Mage Targets for Website

**Description**: Integrate website build with mage tasks (optional).

**Tasks**:
1. Edit `cli/magefile.go`:
   - Add `WebsiteDev()` func: `pnpm --dir ../web run dev`
   - Add `WebsiteBuild()` func: `pnpm --dir ../web run build`
   - Add `WebsitePreview()` func: `pnpm --dir ../web run preview`

2. Test:
   - `mage websiteDev`
   - `mage websiteBuild`

**Acceptance**:
- Mage targets work
- Convenient for developers

---

## IN PROGRESS

(None)

---

## DONE

(Tasks move here as completed)
