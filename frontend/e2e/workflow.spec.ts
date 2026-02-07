import { test, expect } from '@playwright/test';

test.describe('Workflow Editor', () => {
  test('should load the workflow editor page', async ({ page }) => {
    await page.goto('/workflow/new');
    await expect(page.locator('text=Flow')).toBeVisible();
    await expect(page.locator('button:has-text("Add Node")')).toBeVisible();
  });

  test('should be able to add a node', async ({ page }) => {
    await page.goto('/workflow/new');
    
    // Click Add Node button
    await page.click('button:has-text("Add Node")');
    
    // Verify node menu appears
    await expect(page.locator('text=httpRequest')).toBeVisible();
    
    // Click on HTTP Request node type
    await page.click('text=httpRequest');
    
    // Verify node is added to canvas
    await expect(page.locator('[data-id]')).toHaveCount(1);
  });

  test('should navigate to workflow list', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('text=Workflows')).toBeVisible();
  });
});
