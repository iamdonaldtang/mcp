import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Pagination from '@/components/common/Pagination.vue'

describe('Pagination', () => {
  it('displays correct "Showing X-Y of Z" text', () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 20, total: 100 },
    })
    expect(wrapper.text()).toContain('Showing 1–20 of 100')
  })

  it('displays correct range for middle pages', () => {
    const wrapper = mount(Pagination, {
      props: { page: 3, pageSize: 20, total: 100 },
    })
    expect(wrapper.text()).toContain('Showing 41–60 of 100')
  })

  it('caps the end range at total for the last page', () => {
    const wrapper = mount(Pagination, {
      props: { page: 3, pageSize: 20, total: 50 },
    })
    expect(wrapper.text()).toContain('Showing 41–50 of 50')
  })

  it('disables Prev button on first page', () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 20, total: 100 },
    })
    const prevBtn = wrapper.findAll('button').find(b => b.text() === 'Prev')!
    expect(prevBtn.attributes('disabled')).toBeDefined()
  })

  it('disables Next button on last page', () => {
    const wrapper = mount(Pagination, {
      props: { page: 5, pageSize: 20, total: 100 },
    })
    const nextBtn = wrapper.findAll('button').find(b => b.text() === 'Next')!
    expect(nextBtn.attributes('disabled')).toBeDefined()
  })

  it('enables Prev button when not on first page', () => {
    const wrapper = mount(Pagination, {
      props: { page: 2, pageSize: 20, total: 100 },
    })
    const prevBtn = wrapper.findAll('button').find(b => b.text() === 'Prev')!
    expect(prevBtn.attributes('disabled')).toBeUndefined()
  })

  it('enables Next button when not on last page', () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 20, total: 100 },
    })
    const nextBtn = wrapper.findAll('button').find(b => b.text() === 'Next')!
    expect(nextBtn.attributes('disabled')).toBeUndefined()
  })

  it('emits update:page with next page on Next click', async () => {
    const wrapper = mount(Pagination, {
      props: { page: 2, pageSize: 20, total: 100 },
    })
    const nextBtn = wrapper.findAll('button').find(b => b.text() === 'Next')!
    await nextBtn.trigger('click')
    expect(wrapper.emitted('update:page')).toEqual([[3]])
  })

  it('emits update:page with previous page on Prev click', async () => {
    const wrapper = mount(Pagination, {
      props: { page: 3, pageSize: 20, total: 100 },
    })
    const prevBtn = wrapper.findAll('button').find(b => b.text() === 'Prev')!
    await prevBtn.trigger('click')
    expect(wrapper.emitted('update:page')).toEqual([[2]])
  })

  it('emits update:page when a page number button is clicked', async () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 20, total: 60 },
    })
    const pageBtn = wrapper.findAll('button').find(b => b.text() === '2')!
    await pageBtn.trigger('click')
    expect(wrapper.emitted('update:page')).toEqual([[2]])
  })

  it('does not emit when Prev is clicked on first page', async () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 20, total: 100 },
    })
    const prevBtn = wrapper.findAll('button').find(b => b.text() === 'Prev')!
    await prevBtn.trigger('click')
    expect(wrapper.emitted('update:page')).toBeUndefined()
  })

  it('does not emit when Next is clicked on last page', async () => {
    const wrapper = mount(Pagination, {
      props: { page: 5, pageSize: 20, total: 100 },
    })
    const nextBtn = wrapper.findAll('button').find(b => b.text() === 'Next')!
    await nextBtn.trigger('click')
    expect(wrapper.emitted('update:page')).toBeUndefined()
  })

  it('highlights the current page button', () => {
    const wrapper = mount(Pagination, {
      props: { page: 2, pageSize: 20, total: 60 },
    })
    const pageBtn = wrapper.findAll('button').find(b => b.text() === '2')!
    expect(pageBtn.classes()).toContain('bg-community')
  })

  it('computes totalPages correctly', () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 20, total: 45 },
    })
    // 45 / 20 = 2.25 -> ceil = 3 pages
    const pageButtons = wrapper.findAll('button').filter(b => /^\d+$/.test(b.text()))
    expect(pageButtons).toHaveLength(3)
  })

  it('shows ellipsis for many pages', () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 10, total: 200 },
    })
    expect(wrapper.text()).toContain('...')
  })

  it('handles single page total', () => {
    const wrapper = mount(Pagination, {
      props: { page: 1, pageSize: 20, total: 10 },
    })
    const prevBtn = wrapper.findAll('button').find(b => b.text() === 'Prev')!
    const nextBtn = wrapper.findAll('button').find(b => b.text() === 'Next')!
    expect(prevBtn.attributes('disabled')).toBeDefined()
    expect(nextBtn.attributes('disabled')).toBeDefined()
  })
})
