---
title: Website E2E Tests
description: End-to-end tests for azd-exec documentation website
lastUpdated: 2026-01-09
tags: [testing, e2e, playwright, website]
---

# azd-exec Website Tests

Playwright end-to-end tests for the azd-exec documentation website.

## Setup

```bash
# Install dependencies
pnpm install

# Install Playwright browsers
pnpm exec playwright install
```

## Running Tests

```bash
# Run all tests
pnpm test

# Run tests in headed mode (see browser)
pnpm exec playwright test --headed

# Run specific test file
pnpm exec playwright test homepage.spec.ts

# Run tests in debug mode
pnpm exec playwright test --debug

# View test report
pnpm exec playwright show-report
```

## Test Coverage

The test suite covers:

- **Homepage** (`homepage.spec.ts`) - Main landing page functionality, meta tags, dark theme
- **Getting Started** (`getting-started.spec.ts`) - Installation instructions and quick start
- **Examples** (`examples.spec.ts`) - Code examples for bash and PowerShell
- **Security** (`security.spec.ts`) - Security guidelines and best practices
- **CLI Reference** (`cli-reference.spec.ts`) - Command documentation
- **Changelog** (`changelog.spec.ts`) - Version history and release notes
- **Navigation** (`navigation.spec.ts`) - Cross-page navigation and links
- **Accessibility** (`accessibility.spec.ts`) - A11y compliance (WCAG AA)

## Configuration

Tests are configured in `playwright.config.ts` with:
- Multiple browsers (Chromium, Firefox, WebKit)
- Mobile viewports (Pixel 5, iPhone 12)
- Screenshot on failure
- Trace collection on retry
- Automatic dev server startup

## CI/CD

Tests can be run in CI with:

```bash
# Start dev server in background
pnpm dev &

# Run tests
pnpm exec playwright test

# Kill dev server
pkill -f "astro dev"
```

Or let Playwright manage the server automatically (default).

## Writing New Tests

1. Create a new `.spec.ts` file in the `tests/` directory
2. Import Playwright test utilities
3. Use `BASE_URL` constant for consistency
4. Group related tests in `test.describe()` blocks
5. Use semantic locators (`getByRole`, `getByLabel`, etc.)
6. Add assertions with `expect()`

Example:

```typescript
import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('My Feature', () => {
  test('should work correctly', async ({ page }) => {
    await page.goto(`${BASE_URL}/my-page`);
    await expect(page.getByRole('heading')).toBeVisible();
  });
});
```

## Troubleshooting

### Port already in use

If you get a port conflict:

```bash
# Find and kill process on port 4321
netstat -ano | findstr :4321
taskkill /PID <PID> /F
```

### Tests failing locally

1. Clear Astro cache: `rm -rf .astro`
2. Restart dev server
3. Update Playwright browsers: `pnpm exec playwright install`

### Smart quote issues

If you encounter "Unterminated string literal" errors in Layout.astro, check for smart quotes (') and replace with straight quotes (').

```bash
# PowerShell: Replace smart quotes
$content = Get-Content "src/components/Layout.astro" -Raw
$fixed = $content.Replace([char]0x2019, [char]0x0027)
$fixed | Out-File "src/components/Layout.astro" -Encoding UTF8 -NoNewline
```
