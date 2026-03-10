import { test, expect } from '@playwright/test'

test.describe('Community Creation Wizard', () => {
  test.beforeEach(async ({ page }) => {
    // Login first
    await page.goto('/login')
    await page.fill('input[type="email"]', 'test@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button:has-text("Login"), button:has-text("Sign In")')
    await page.waitForURL('**/b/dashboard', { timeout: 10000 })
  })

  test('navigate to community empty state', async ({ page }) => {
    // Navigate to community via sidebar or product card
    await page.click('text=Community, a[href*="community"]')
    await expect(page.locator('text=Welcome to Community')).toBeVisible()
  })

  test('start wizard with strategy template', async ({ page }) => {
    await page.goto('/b/community/empty')
    // Click a strategy card then CTA
    await page.click('text=Create Your First Community, button:has-text("Get Started")')
    await page.waitForURL('**/wizard/step-1**')
    // Step 1: Fill name and description
    await page.fill('input[placeholder*="name" i], input:first-of-type', 'My Test Community')
    // Verify stepper shows step 1 active
    await expect(page.locator('text=Customize')).toBeVisible()
  })

  test('wizard step 1 validation', async ({ page }) => {
    await page.goto('/b/community/wizard/step-1')
    // Type short name
    await page.fill('input:first-of-type', 'ab')
    // Next button should be disabled
    const nextBtn = page.locator('button:has-text("Next")')
    await expect(nextBtn).toBeDisabled()
    // Type valid name
    await page.fill('input:first-of-type', 'My Community')
    // Still need description
    await expect(nextBtn).toBeDisabled()
  })
})
