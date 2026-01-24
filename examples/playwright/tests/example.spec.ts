import { test, expect } from '@playwright/test';

test('page has title', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveTitle(/Example Domain/);
});

test('page has content', async ({ page }) => {
  await page.goto('/');
  await expect(page.locator('h1')).toContainText('Example Domain');
});
