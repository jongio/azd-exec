import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('Homepage', () => {
  test('should load successfully', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    await expect(page).toHaveTitle(/azd exec/i);
  });

  test('should have correct heading', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    const heading = page.getByRole('heading', { level: 1 });
    await expect(heading).toBeVisible();
  });

  test('should display install instructions', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    // Check for install code blocks or tabs
    const installSection = page.locator('text=/install/i').first();
    await expect(installSection).toBeVisible();
  });

  test('should have navigation links', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    // Check for header navigation
    const nav = page.locator('nav,header').first();
    await expect(nav).toBeVisible();
  });

  test('should have footer', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    const footer = page.locator('footer');
    await expect(footer).toBeVisible();
  });

  test('should use dark theme by default', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    const html = page.locator('html');
    await expect(html).toHaveAttribute('data-theme', 'dark');
  });

  test('should have proper meta tags', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    // Check for description meta tag
    const description = page.locator('meta[name="description"]');
    await expect(description).toHaveAttribute('content', /.+/);
    
    // Check for Open Graph tags
    const ogTitle = page.locator('meta[property="og:title"]');
    await expect(ogTitle).toHaveAttribute('content', /.+/);
  });

  test('should be responsive', async ({ page }) => {
    // Test mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(`${BASE_URL}/`);
    await expect(page.locator('body')).toBeVisible();
    
    // Test tablet viewport
    await page.setViewportSize({ width: 768, height: 1024 });
    await expect(page.locator('body')).toBeVisible();
    
    // Test desktop viewport
    await page.setViewportSize({ width: 1920, height: 1080 });
    await expect(page.locator('body')).toBeVisible();
  });
});
