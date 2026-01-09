---
title: Website Testing Setup
description: Playwright test setup for azd-exec documentation website
lastUpdated: 2026-01-09
tags: [testing, playwright, e2e, website]
---

# Website Testing Setup

## Quick Start

```bash
cd web

# Install dependencies (includes Playwright)
pnpm install

# Install Playwright browsers
pnpm exec playwright install

# Run tests
pnpm test
```

## What Was Created

### Test Files

8 comprehensive test suites covering all pages:

1. **homepage.spec.ts** - Homepage tests (dark theme, meta tags, responsive)
2. **getting-started.spec.ts** - Installation and quick start  
3. **examples.spec.ts** - Code examples (bash, PowerShell)
4. **security.spec.ts** - Security guidelines
5. **cli-reference.spec.ts** - CLI command documentation
6. **changelog.spec.ts** - Version history
7. **navigation.spec.ts** - Cross-page navigation
8. **accessibility.spec.ts** - WCAG AA compliance

### Configuration

- **playwright.config.ts** - Multi-browser (Chrome, Firefox, Safari) + mobile testing
- **package.json** - Updated with test scripts and Playwright dependency

### Features

- ✅ Dark theme verification
- ✅ Responsive design testing (mobile, tablet, desktop)
- ✅ SEO meta tags validation
- ✅ Accessibility (A11y) compliance
- ✅ Cross-browser testing
- ✅ Screenshot on failure
- ✅ Automatic dev server management

## Known Issue

The Layout.astro file has a smart quote character (') instead of a straight quote (') on line 19, causing a build error. This needs to be fixed before tests can run successfully.

### Fix

```powershell
# Replace smart quote with straight quote
$content = Get-Content "src/components/Layout.astro" -Raw -Encoding UTF8
$fixed = $content.Replace([char]0x2019, [char]0x0027)
$fixed | Out-File "src/components/Layout.astro" -Encoding UTF8 -NoNewline
```

Or manually edit line 19 and replace the apostrophe in "your Azure Developer CLI" with a straight quote.

## Running Tests

```bash
# All tests
pnpm test

# With browser UI
pnpm test:headed

# Debug mode
pnpm test:debug

# View last report
pnpm test:report

# Specific file
pnpm exec playwright test homepage.spec.ts

# Specific browser
pnpm exec playwright test --project=chromium
```

## CI/CD Integration

Tests automatically start/stop the dev server. For CI pipelines, set `CI=true`:

```bash
CI=true pnpm test
```

## Next Steps

1. Fix the Layout.astro smart quote issue
2. Run `pnpm install` to add Playwright
3. Run `pnpm exec playwright install` for browsers
4. Run `pnpm test` to execute all tests
5. Review test report with `pnpm test:report`

See `tests/README.md` for detailed documentation.
