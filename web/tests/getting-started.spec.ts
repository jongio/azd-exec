import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('Getting Started Page', () => {
  test('should load successfully', async ({ page }) => {
    await page.goto(`${BASE_URL}/getting-started`);
    await expect(page).toHaveTitle(/getting started.*azd exec/i);
  });

  test('should have installation section', async ({ page }) => {
    await page.goto(`${BASE_URL}/getting-started`);
    const installHeading = page.getByRole('heading', { name: /^install azd exec$/i });
    await expect(installHeading).toBeVisible();
  });

  test('should display code examples', async ({ page }) => {
    await page.goto(`${BASE_URL}/getting-started`);
    // Check for code content (expressive-code uses figure > pre > code structure)
    const codeBlocks = page.locator('pre code, .expressive-code code, [class*="code"]');
    await expect(codeBlocks.first()).toBeAttached();
    // Verify there's actual code content
    const codeContent = await codeBlocks.first().textContent();
    expect(codeContent?.trim()).not.toBe('');
  });

  test('should have quick start guide', async ({ page }) => {
    await page.goto(`${BASE_URL}/getting-started`);
    const quickStart = page.locator('text=/quick start|usage|how to use/i').first();
    await expect(quickStart).toBeVisible();
  });

  test('should contain executable commands', async ({ page }) => {
    await page.goto(`${BASE_URL}/getting-started`);
    // Look for azd exec commands
    const azdExecCommand = page.locator('text=/azd exec/i').first();
    await expect(azdExecCommand).toBeVisible();
  });
});
