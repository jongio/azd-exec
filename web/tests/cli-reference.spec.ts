import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('CLI Reference', () => {
  test('should load CLI reference index', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/cli`);
    await expect(page).toHaveTitle(/cli.*reference|reference.*cli/i);
  });

  test('should list available commands', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/cli`);
    const commandsList = page.locator('text=/command|usage/i').first();
    await expect(commandsList).toBeVisible();
  });

  test('should have version command documentation', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/cli/version`);
    await expect(page).toHaveTitle(/version.*azd exec/i);
  });

  test('should show command syntax', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/cli/version`);
    const syntax = page.locator('pre code, .syntax, .usage').first();
    await expect(syntax).toBeVisible();
  });

  test('should document command flags', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/cli/version`);
    const flags = page.locator('text=/flag|option|argument|-/').first();
    await expect(flags).toBeVisible();
  });
});
