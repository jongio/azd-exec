import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('Navigation', () => {
  test('should navigate from home to getting started', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    await page.click('text=/getting started|get started/i');
    await expect(page).toHaveURL(/getting-started/);
  });

  test('should navigate from home to examples', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    await page.click('text=/example/i');
    await expect(page).toHaveURL(/examples/);
  });

  test('should navigate to security page', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    const securityLink = page.locator('a[href*="security"]').first();
    await securityLink.click();
    await expect(page).toHaveURL(/security/);
  });

  test('should navigate to CLI reference', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    const cliLink = page.locator('a[href*="reference"]').first();
    await cliLink.click();
    await expect(page).toHaveURL(/reference/);
  });

  test('should have working logo/home link', async ({ page }) => {
    await page.goto(`${BASE_URL}/examples`);
    const homeLink = page.locator('a[href*="/"]').first();
    await homeLink.click();
    await expect(page).toHaveURL(new RegExp(`${BASE_URL}/?$`));
  });
});
