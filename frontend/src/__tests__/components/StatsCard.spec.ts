import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatsCard from '@/components/common/StatsCard.vue'

describe('StatsCard', () => {
  const baseProps = {
    label: 'Total Users',
    value: 12500,
    icon: 'group',
  }

  it('renders the label text', () => {
    const wrapper = mount(StatsCard, { props: baseProps })
    expect(wrapper.text()).toContain('Total Users')
  })

  it('formats numeric values with toLocaleString', () => {
    const wrapper = mount(StatsCard, { props: baseProps })
    // 12500 -> "12,500" (en locale)
    expect(wrapper.text()).toContain('12,500')
  })

  it('renders string values as-is without formatting', () => {
    const wrapper = mount(StatsCard, {
      props: { ...baseProps, value: '$1,234.56' },
    })
    expect(wrapper.text()).toContain('$1,234.56')
  })

  it('renders the icon text', () => {
    const wrapper = mount(StatsCard, { props: baseProps })
    const iconSpan = wrapper.find('.material-symbols-rounded')
    expect(iconSpan.text()).toBe('group')
  })

  it('applies custom iconColor when provided', () => {
    const wrapper = mount(StatsCard, {
      props: { ...baseProps, iconColor: '#48BB78' },
    })
    const iconSpan = wrapper.find('.material-symbols-rounded')
    expect(iconSpan.attributes('style')).toContain('color: #48BB78')
  })

  it('uses default icon color #94A3B8 when iconColor is not provided', () => {
    const wrapper = mount(StatsCard, { props: baseProps })
    const iconSpan = wrapper.find('.material-symbols-rounded')
    expect(iconSpan.attributes('style')).toContain('color: #94A3B8')
  })

  it('renders trend text when provided', () => {
    const wrapper = mount(StatsCard, {
      props: { ...baseProps, trend: '+12% this week' },
    })
    expect(wrapper.text()).toContain('+12% this week')
  })

  it('does not render trend element when trend is not provided', () => {
    const wrapper = mount(StatsCard, { props: baseProps })
    // The trend div uses v-if, so it should not exist
    const allText = wrapper.text()
    // Just ensure no extra div beyond label, value, icon
    expect(allText).toBe('Total Usersgroup12,500')
  })

  it('handles zero value correctly', () => {
    const wrapper = mount(StatsCard, {
      props: { ...baseProps, value: 0 },
    })
    expect(wrapper.text()).toContain('0')
  })

  it('handles large numeric values', () => {
    const wrapper = mount(StatsCard, {
      props: { ...baseProps, value: 1000000 },
    })
    expect(wrapper.text()).toContain('1,000,000')
  })
})
