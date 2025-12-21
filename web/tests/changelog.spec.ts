import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('Changelog', () => {
  test('should load changelog page', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/changelog`);
    await expect(page).toHaveTitle(/changelog.*azd exec/i);
  });

  test('should display version history', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/changelog`);
    // Look for version numbers (semver format)
    const version = page.locator('text=/v?\\d+\\.\\d+\\.\\d+|version/i').first();
    await expect(version).toBeVisible();
  });

  test('should show release dates', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/changelog`);
    const date = page.locator('text=/\\d{4}-\\d{2}-\\d{2}|\\d{1,2}\\/\\d{1,2}\\/\\d{4}|january|february|march|april|may|june|july|august|september|october|november|december/i').first();
    await expect(date).toBeVisible();
  });

  test('should list changes', async ({ page }) => {
    await page.goto(`${BASE_URL}/reference/changelog`);
    const changes = page.locator('ul li, ol li').first();
    await expect(changes).toBeVisible();
  });
});
