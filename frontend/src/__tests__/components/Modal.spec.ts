import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Modal from '@/components/common/Modal.vue'

describe('Modal', () => {
  const baseProps = {
    open: true,
    title: 'Confirm Action',
  }

  it('renders content when open is true', () => {
    const wrapper = mount(Modal, {
      props: baseProps,
      slots: { default: '<p>Modal body</p>' },
    })
    expect(wrapper.text()).toContain('Confirm Action')
    expect(wrapper.text()).toContain('Modal body')
  })

  it('does not render content when open is false', () => {
    const wrapper = mount(Modal, {
      props: { ...baseProps, open: false },
      slots: { default: '<p>Modal body</p>' },
    })
    expect(wrapper.text()).not.toContain('Confirm Action')
    expect(wrapper.text()).not.toContain('Modal body')
  })

  it('renders the title in the header', () => {
    const wrapper = mount(Modal, { props: baseProps })
    expect(wrapper.find('h2').text()).toBe('Confirm Action')
  })

  it('emits close when backdrop is clicked', async () => {
    const wrapper = mount(Modal, { props: baseProps })
    // Backdrop is the div with bg-black/60
    const backdrop = wrapper.find('.bg-black\\/60')
    await backdrop.trigger('click')
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('emits close when close button is clicked', async () => {
    const wrapper = mount(Modal, { props: baseProps })
    // Close button contains the "close" icon text
    const closeBtn = wrapper.find('button')
    await closeBtn.trigger('click')
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('applies default maxWidth of 480px', () => {
    const wrapper = mount(Modal, { props: baseProps })
    const content = wrapper.find('.relative')
    expect(content.attributes('style')).toContain('max-width: 480px')
  })

  it('applies custom maxWidth when provided', () => {
    const wrapper = mount(Modal, {
      props: { ...baseProps, maxWidth: '640px' },
    })
    const content = wrapper.find('.relative')
    expect(content.attributes('style')).toContain('max-width: 640px')
  })

  it('renders default slot in the body', () => {
    const wrapper = mount(Modal, {
      props: baseProps,
      slots: { default: '<div class="body-content">Hello</div>' },
    })
    expect(wrapper.find('.body-content').exists()).toBe(true)
    expect(wrapper.find('.body-content').text()).toBe('Hello')
  })

  it('renders footer slot when provided', () => {
    const wrapper = mount(Modal, {
      props: baseProps,
      slots: {
        default: '<p>Body</p>',
        footer: '<button>Save</button><button>Cancel</button>',
      },
    })
    expect(wrapper.text()).toContain('Save')
    expect(wrapper.text()).toContain('Cancel')
  })

  it('does not render footer section when footer slot is empty', () => {
    const wrapper = mount(Modal, {
      props: baseProps,
      slots: { default: '<p>Body only</p>' },
    })
    // Footer div has border-t border-border class; count bordered sections
    const footerDiv = wrapper.findAll('.border-t').filter(el => el.classes().includes('flex'))
    // The footer should not exist when no footer slot is provided
    // We check by looking for the justify-end flex container
    const justifyEnd = wrapper.findAll('.justify-end')
    expect(justifyEnd).toHaveLength(0)
  })

  it('toggles visibility reactively when open prop changes', async () => {
    const wrapper = mount(Modal, {
      props: { ...baseProps, open: false },
      slots: { default: '<p>Content</p>' },
    })
    expect(wrapper.text()).not.toContain('Content')

    await wrapper.setProps({ open: true })
    expect(wrapper.text()).toContain('Content')

    await wrapper.setProps({ open: false })
    expect(wrapper.text()).not.toContain('Content')
  })
})
