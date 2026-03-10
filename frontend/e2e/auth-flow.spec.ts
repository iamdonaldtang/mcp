import { test, expect } from '@playwright/test'

test.describe('B-End Auth Flow', () => {
  test('register new account', async ({ page }) => {
    await page.goto('/login')
    await page.fill('[data-testid="email-input"], input[type="email"]', 'test@example.com')
    await page.fill('[data-testid="password-input"], input[type="password"]', 'password123')
    await page.fill('[data-testid="name-input"], input[placeholder*="name" i]', 'Test User')
    // Look for register/signup button
    await page.click('button:has-text("Register"), button:has-text("Sign Up")')
    // Should redirect to dashboard
    await page.waitForURL('**/b/dashboard', { timeout: 10000 })
    await expect(page.locator('text=Welcome')).toBeVisible()
  })

  test('login with existing account', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'test@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button:has-text("Login"), button:has-text("Sign In")')
    await page.waitForURL('**/b/dashboard', { timeout: 10000 })
  })

  test('access protected route without auth redirects to login', async ({ page }) => {
    await page.goto('/b/dashboard')
    await page.waitForURL('**/login', { timeout: 10000 })
  })
})
