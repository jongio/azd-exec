import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('Examples Page', () => {
  test('should load successfully', async ({ page }) => {
    await page.goto(`${BASE_URL}/examples`);
    await expect(page).toHaveTitle(/examples.*azd exec/i);
  });

  test('should display example scripts', async ({ page }) => {
    await page.goto(`${BASE_URL}/examples`);
    const examples = page.locator('pre code, .example, .code-block');
    await expect(examples.first()).toBeVisible();
  });

  test('should have multiple example categories', async ({ page }) => {
    await page.goto(`${BASE_URL}/examples`);
    // Count headings which typically separate example categories
    const headings = page.locator('h2, h3');
    const count = await headings.count();
    expect(count).toBeGreaterThan(1);
  });

  test('should show bash examples', async ({ page }) => {
    await page.goto(`${BASE_URL}/examples`);
    const bashExample = page.locator('text=/bash|sh|\.sh/i').first();
    await expect(bashExample).toBeVisible();
  });

  test('should show PowerShell examples', async ({ page }) => {
    await page.goto(`${BASE_URL}/examples`);
    const pwshExample = page.locator('text=/powershell|ps1|\.ps1/i').first();
    await expect(pwshExample).toBeVisible();
  });
});
