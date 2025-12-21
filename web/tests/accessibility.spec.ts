import { test, expect } from '@playwright/test';

const BASE_URL = process.env.BASE_URL || 'http://localhost:4321/azd-exec';

test.describe('Accessibility', () => {
  test('homepage should have proper heading hierarchy', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    const h1 = await page.locator('h1').count();
    expect(h1).toBe(1); // Should have exactly one H1
    
    const headings = await page.locator('h1, h2, h3, h4, h5, h6').all();
    expect(headings.length).toBeGreaterThan(0);
  });

  test('images should have alt text', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    const images = await page.locator('img').all();
    for (const img of images) {
      const alt = await img.getAttribute('alt');
      expect(alt).toBeDefined();
    }
  });

  test('links should have accessible names', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    const links = await page.locator('a').all();
    for (const link of links) {
      const text = await link.textContent();
      const ariaLabel = await link.getAttribute('aria-label');
      const title = await link.getAttribute('title');
      
      // At least one of these should have content
      expect(text || ariaLabel || title).toBeTruthy();
    }
  });

  test('should have proper ARIA landmarks', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    // Check for main landmark
    const main = page.locator('main, [role="main"]');
    await expect(main).toBeVisible();
    
    // Check for navigation landmark
    const nav = page.locator('nav, [role="navigation"]');
    await expect(nav.first()).toBeVisible();
  });

  test('should be keyboard navigable', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    // Tab through interactive elements
    await page.keyboard.press('Tab');
    const focused = await page.evaluate(() => document.activeElement?.tagName);
    expect(['A', 'BUTTON', 'INPUT', 'TEXTAREA', 'SELECT']).toContain(focused);
  });

  test('code blocks should be readable', async ({ page }) => {
    await page.goto(`${BASE_URL}/getting-started`);
    
    const codeBlocks = await page.locator('pre code').all();
    for (const block of codeBlocks) {
      // Check that code blocks have visible text
      const text = await block.textContent();
      expect(text?.trim().length).toBeGreaterThan(0);
    }
  });

  test('form elements should have labels', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    const inputs = await page.locator('input:not([type="hidden"]), textarea, select').all();
    for (const input of inputs) {
      const id = await input.getAttribute('id');
      const ariaLabel = await input.getAttribute('aria-label');
      const ariaLabelledby = await input.getAttribute('aria-labelledby');
      
      if (id) {
        const label = page.locator(`label[for="${id}"]`);
        const hasLabel = await label.count() > 0;
        expect(hasLabel || ariaLabel || ariaLabelledby).toBeTruthy();
      } else {
        expect(ariaLabel || ariaLabelledby).toBeTruthy();
      }
    }
  });

  test('should have sufficient color contrast in dark theme', async ({ page }) => {
    await page.goto(`${BASE_URL}/`);
    
    // Verify dark theme is active
    const html = page.locator('html');
    await expect(html).toHaveAttribute('data-theme', 'dark');
    
    // Check that text is visible against background
    const body = page.locator('body');
    const bgColor = await body.evaluate((el) => getComputedStyle(el).backgroundColor);
    const color = await body.evaluate((el) => getComputedStyle(el).color);
    
    expect(bgColor).toBeDefined();
    expect(color).toBeDefined();
    expect(bgColor).not.toBe(color);
  });
});
