import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EmptyState from '@/components/common/EmptyState.vue'

describe('EmptyState', () => {
  const baseProps = {
    icon: 'inbox',
    title: 'No items yet',
    description: 'Create your first item to get started.',
  }

  it('renders the icon', () => {
    const wrapper = mount(EmptyState, { props: baseProps })
    const iconSpan = wrapper.find('.material-symbols-rounded')
    expect(iconSpan.text()).toBe('inbox')
  })

  it('renders the title', () => {
    const wrapper = mount(EmptyState, { props: baseProps })
    expect(wrapper.find('h3').text()).toBe('No items yet')
  })

  it('renders the description', () => {
    const wrapper = mount(EmptyState, { props: baseProps })
    expect(wrapper.find('p').text()).toBe('Create your first item to get started.')
  })

  it('renders CTA button when actionLabel and actionRoute are provided', () => {
    const wrapper = mount(EmptyState, {
      props: {
        ...baseProps,
        actionLabel: 'Create Item',
        actionRoute: '/items/new',
      },
    })
    // router-link is stubbed as <a :href="to">
    const link = wrapper.find('a')
    expect(link.exists()).toBe(true)
    expect(link.text()).toBe('Create Item')
    expect(link.attributes('href')).toBe('/items/new')
  })

  it('does not render CTA button when actionLabel is missing', () => {
    const wrapper = mount(EmptyState, {
      props: { ...baseProps, actionRoute: '/items/new' },
    })
    const link = wrapper.find('a')
    expect(link.exists()).toBe(false)
  })

  it('does not render CTA button when actionRoute is missing', () => {
    const wrapper = mount(EmptyState, {
      props: { ...baseProps, actionLabel: 'Create Item' },
    })
    const link = wrapper.find('a')
    expect(link.exists()).toBe(false)
  })

  it('does not render CTA button when neither action prop is provided', () => {
    const wrapper = mount(EmptyState, { props: baseProps })
    const link = wrapper.find('a')
    expect(link.exists()).toBe(false)
  })

  it('renders default slot content', () => {
    const wrapper = mount(EmptyState, {
      props: baseProps,
      slots: {
        default: '<div class="custom-slot">Custom content</div>',
      },
    })
    expect(wrapper.find('.custom-slot').exists()).toBe(true)
    expect(wrapper.find('.custom-slot').text()).toBe('Custom content')
  })

  it('renders both CTA and slot content together', () => {
    const wrapper = mount(EmptyState, {
      props: {
        ...baseProps,
        actionLabel: 'Create',
        actionRoute: '/create',
      },
      slots: {
        default: '<p class="extra">Extra info</p>',
      },
    })
    expect(wrapper.find('a').exists()).toBe(true)
    expect(wrapper.find('.extra').exists()).toBe(true)
  })
})
