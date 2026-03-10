import { createPinia, setActivePinia } from 'pinia'
import { mount, type ComponentMountingOptions } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import type { Component } from 'vue'

/**
 * Create a fresh Pinia instance for testing
 */
export function createTestPinia() {
  const pinia = createPinia()
  setActivePinia(pinia)
  return pinia
}

/**
 * Create a test router with optional routes
 */
export function createTestRouter(routes: { path: string; name?: string; component: Component }[] = []) {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
      ...routes,
    ],
  })
}

/**
 * Mount a component with Pinia and Router pre-configured
 */
export function mountWithPlugins<T extends Component>(
  component: T,
  options: ComponentMountingOptions<T> & { routes?: { path: string; name?: string; component: Component }[] } = {},
) {
  const pinia = createTestPinia()
  const { routes, ...mountOptions } = options
  const router = createTestRouter(routes)

  return mount(component, {
    ...mountOptions,
    global: {
      ...mountOptions.global,
      plugins: [pinia, router, ...(mountOptions.global?.plugins ?? [])],
      stubs: {
        Teleport: { template: '<div><slot /></div>' },
        ...mountOptions.global?.stubs,
      },
    },
  })
}

/**
 * Flush all pending promises (useful for async component updates)
 */
export async function flushPromises() {
  await new Promise(resolve => setTimeout(resolve, 0))
}

/**
 * Create a mock Axios response
 */
export function mockApiResponse<T>(data: T, status = 200) {
  return Promise.resolve({
    data: { code: 0, data, message: 'success' },
    status,
    statusText: 'OK',
    headers: {},
    config: {} as never,
  })
}

/**
 * Create a mock Axios error response
 */
export function mockApiError(status: number, message = 'error') {
  const error = new Error(message) as Error & { response: { status: number; data: unknown } }
  error.response = { status, data: { code: status, data: null, message } }
  return Promise.reject(error)
}
