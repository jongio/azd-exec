import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('Security Page', () => {
  test('should load successfully', async ({ page }) => {
    await page.goto(`${BASE_URL}/security`);
    await expect(page).toHaveTitle(/security.*azd exec/i);
  });

  test('should explain security considerations', async ({ page }) => {
    await page.goto(`${BASE_URL}/security`);
    const securityContent = page.locator('text=/security|credentials|authentication|permissions/i').first();
    await expect(securityContent).toBeVisible();
  });

  test('should mention Azure credentials', async ({ page }) => {
    await page.goto(`${BASE_URL}/security`);
    const azureCredentials = page.locator('text=/azure.*credential|credential.*azure/i').first();
    await expect(azureCredentials).toBeVisible();
  });

  test('should discuss environment variables', async ({ page }) => {
    await page.goto(`${BASE_URL}/security`);
    const envVars = page.locator('text=/environment variable|env var/i').first();
    await expect(envVars).toBeVisible();
  });

  test('should have best practices section', async ({ page }) => {
    await page.goto(`${BASE_URL}/security`);
    const bestPractices = page.locator('text=/best practice|recommendation|guideline/i').first();
    await expect(bestPractices).toBeVisible();
  });
});
