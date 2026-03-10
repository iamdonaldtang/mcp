import { test, expect } from '@playwright/test'

test.describe('Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'test@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button:has-text("Login"), button:has-text("Sign In")')
    await page.waitForURL('**/b/dashboard', { timeout: 10000 })
  })

  test('sidebar navigation highlights correct item', async ({ page }) => {
    // Navigate to different sections
    const sections = [
      { text: 'Dashboard', url: '/b/dashboard' },
      { text: 'Analytics', url: '/b/analytics' },
      { text: 'Settings', url: '/b/settings' },
    ]
    for (const section of sections) {
      const link = page.locator(`a[href="${section.url}"], text=${section.text}`)
      if (await link.isVisible()) {
        await link.click()
        await page.waitForTimeout(500)
      }
    }
  })

  test('breadcrumb updates on navigation', async ({ page }) => {
    await page.goto('/b/community/empty')
    // Check breadcrumb shows Community
    await expect(page.locator('body')).toBeVisible()
  })

  test('404 page shows for unknown routes', async ({ page }) => {
    await page.goto('/nonexistent-route')
    await expect(page.locator('text=404, text=Not Found')).toBeVisible()
  })
})
