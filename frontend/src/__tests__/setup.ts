import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {}
  return {
    getItem: vi.fn((key: string) => store[key] ?? null),
    setItem: vi.fn((key: string, value: string) => { store[key] = value }),
    removeItem: vi.fn((key: string) => { delete store[key] }),
    clear: vi.fn(() => { store = {} }),
    get length() { return Object.keys(store).length },
    key: vi.fn((index: number) => Object.keys(store)[index] ?? null),
  }
})()
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock window.location
Object.defineProperty(window, 'location', {
  value: { href: '', pathname: '/', search: '', hash: '' },
  writable: true,
})

// Global stubs for router-link and router-view
config.global.stubs = {
  'router-link': {
    template: '<a :href="to"><slot /></a>',
    props: ['to'],
  },
  'router-view': { template: '<div />' },
  Teleport: { template: '<div><slot /></div>' },
}

// Reset mocks between tests
afterEach(() => {
  vi.restoreAllMocks()
  localStorageMock.clear()
})
