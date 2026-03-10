import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Cross-page store cascade refresh mechanism.
 *
 * Each store domain (community, whitelabel, c-end) tracks a lastModified timestamp.
 * Pages call `checkAndRefresh(fetcher)` on mount — if the timestamp has changed
 * since their last fetch, they re-fetch data.
 *
 * Usage in stores:
 *   const { lastModified, notifyChange } = useCascadeRefresh('community')
 *   // Call notifyChange() after any mutation (create/update/delete)
 *
 * Usage in pages:
 *   const { checkAndRefresh } = useCascadeRefresh('community')
 *   onMounted(() => checkAndRefresh(fetchData))
 */

const domainTimestamps: Record<string, number> = {}
const listeners: Record<string, Set<() => void>> = {}

export function useCascadeRefresh(domain: string) {
  const lastFetched = ref(0)

  /**
   * Notify all listeners that data in this domain has changed
   */
  function notifyChange() {
    domainTimestamps[domain] = Date.now()
    if (listeners[domain]) {
      listeners[domain].forEach(fn => fn())
    }
  }

  /**
   * Check if the domain has been modified since last fetch; if so, call the fetcher
   */
  async function checkAndRefresh(fetcher: () => Promise<void>) {
    const domainTs = domainTimestamps[domain] ?? 0
    if (domainTs > lastFetched.value) {
      await fetcher()
      lastFetched.value = Date.now()
    } else if (lastFetched.value === 0) {
      // First mount — always fetch
      await fetcher()
      lastFetched.value = Date.now()
    }
  }

  /**
   * Subscribe to changes (auto-cleanup on unmount)
   */
  function onDomainChange(callback: () => void) {
    if (!listeners[domain]) {
      listeners[domain] = new Set()
    }
    listeners[domain].add(callback)

    onUnmounted(() => {
      listeners[domain]?.delete(callback)
    })
  }

  return {
    lastModified: domainTimestamps,
    notifyChange,
    checkAndRefresh,
    onDomainChange,
  }
}

/**
 * Reset all timestamps (useful for testing)
 */
export function resetCascadeRefresh() {
  Object.keys(domainTimestamps).forEach(k => delete domainTimestamps[k])
  Object.keys(listeners).forEach(k => delete listeners[k])
}
