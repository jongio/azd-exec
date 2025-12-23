import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

// Helper to handle navigation on both desktop and mobile
async function navigateToLink(page: any, linkText: string) {
  const viewportWidth = page.viewportSize()?.width || 1280;
  
  if (viewportWidth < 768) {
    // Mobile: open menu first, then find link in mobile nav
    const menuButton = page.locator('#mobile-menu-toggle');
    if (await menuButton.isVisible()) {
      await menuButton.click();
      await page.waitForTimeout(500); // Wait for menu animation
    }
    // Scope to mobile nav
    await page.locator('#mobile-menu').locator(linkText).first().click();
  } else {
    // Desktop: scope to desktop nav
    await page.locator('.nav-desktop').locator(linkText).first().click();
  }
}

test.describe('Navigation', () => {
  test('should navigate from home to getting started', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    await navigateToLink(page, 'text=/getting started|get started/i');
    await expect(page).toHaveURL(/getting-started/);
  });

  test('should navigate from home to examples', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    await navigateToLink(page, 'text=/example/i');
    await expect(page).toHaveURL(/examples/);
  });

  test('should navigate to security page', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    // This test doesn't use the helper because security link is in page content
    const securityLink = page.locator('a[href*="security"]').first();
    await securityLink.click({ force: true });
    await expect(page).toHaveURL(/security/);
  });

  test('should navigate to CLI reference', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    await navigateToLink(page, 'a[href*="reference"]');
    await expect(page).toHaveURL(/reference/);
  });

  test('should have working logo/home link', async ({ page }) => {
    await page.goto(`${BASE_URL}/examples`);
    const homeLink = page.locator('a[href*="/"]').first();
    await homeLink.click();
    await expect(page).toHaveURL(new RegExp(`${BASE_URL}/?$`));
  });
});
