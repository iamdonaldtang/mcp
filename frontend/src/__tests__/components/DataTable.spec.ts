import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DataTable from '@/components/common/DataTable.vue'

const columns = [
  { key: 'name', label: 'Name', width: '200px' },
  { key: 'status', label: 'Status' },
  { key: 'date', label: 'Date', width: '150px' },
]

describe('DataTable', () => {
  it('renders all column headers', () => {
    const wrapper = mount(DataTable, { props: { columns } })
    const headers = wrapper.findAll('th')
    expect(headers).toHaveLength(3)
    expect(headers[0].text()).toBe('Name')
    expect(headers[1].text()).toBe('Status')
    expect(headers[2].text()).toBe('Date')
  })

  it('applies width style to columns that specify it', () => {
    const wrapper = mount(DataTable, { props: { columns } })
    const headers = wrapper.findAll('th')
    expect(headers[0].attributes('style')).toContain('width: 200px')
    expect(headers[1].attributes('style')).toBeUndefined()
    expect(headers[2].attributes('style')).toContain('width: 150px')
  })

  it('renders default slot content when not loading', () => {
    const wrapper = mount(DataTable, {
      props: { columns },
      slots: {
        default: '<tr><td>Row 1</td></tr><tr><td>Row 2</td></tr>',
      },
    })
    const rows = wrapper.findAll('tbody tr')
    expect(rows).toHaveLength(2)
    expect(rows[0].text()).toBe('Row 1')
  })

  it('shows loading spinner when loading is true', () => {
    const wrapper = mount(DataTable, {
      props: { columns, loading: true },
      slots: {
        default: '<tr><td>Should not appear</td></tr>',
      },
    })
    const tbody = wrapper.find('tbody')
    expect(tbody.text()).toContain('progress_activity')
    expect(tbody.text()).not.toContain('Should not appear')
  })

  it('sets colspan to match number of columns when loading', () => {
    const wrapper = mount(DataTable, {
      props: { columns, loading: true },
    })
    const td = wrapper.find('tbody td')
    expect(td.attributes('colspan')).toBe('3')
  })

  it('does not show loading spinner when loading is false', () => {
    const wrapper = mount(DataTable, {
      props: { columns, loading: false },
    })
    expect(wrapper.find('tbody').text()).not.toContain('progress_activity')
  })

  it('renders with no rows when slot is empty and not loading', () => {
    const wrapper = mount(DataTable, { props: { columns } })
    const rows = wrapper.findAll('tbody tr')
    expect(rows).toHaveLength(0)
  })

  it('renders table within a wrapper div', () => {
    const wrapper = mount(DataTable, { props: { columns } })
    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.element.tagName).toBe('DIV')
  })
})
