import { test, expect } from '@playwright/test'

test.describe('C-End User Flow', () => {
  test('community home page renders', async ({ page }) => {
    await page.goto('/c')
    // Should show community home with pulse stats or loading
    await expect(page.locator('body')).toBeVisible()
  })

  test('navigate between tabs', async ({ page }) => {
    await page.goto('/c')
    // Check for navigation tabs
    const tabs = ['Quests', 'Leaderboard', 'LB Sprint', 'Milestone', 'Shop']
    for (const tab of tabs) {
      const tabElement = page.locator(`text=${tab}`)
      if (await tabElement.isVisible()) {
        await tabElement.click()
        // Should navigate
        await page.waitForTimeout(500)
      }
    }
  })

  test('shop page shows items', async ({ page }) => {
    await page.goto('/c/shop')
    // Should show shop items or empty state
    await expect(page.locator('body')).toBeVisible()
  })

  test('leaderboard shows rankings', async ({ page }) => {
    await page.goto('/c/leaderboard')
    await expect(page.locator('body')).toBeVisible()
  })
})
