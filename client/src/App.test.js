import React from 'react';
import { shallow, mount } from 'enzyme';
import renderer from 'react-test-renderer'
import { MemoryRouter } from 'react-router'
import App from './App';

describe('App', () => {
  it('renders without crashing', () => {
    const wrapper = shallow(<App />)
    expect(wrapper.exists()).toBe(true)
  });

  it('renders the logo', () => {
    const wrapper = shallow(<App />)
    expect(wrapper.exists()).toBe(true)
  })
})

describe('routing', () => {
  const setup = (path) => {
    const wrapper = mount(
      <MemoryRouter initialEntries={[path]}>
        <App />
      </MemoryRouter>
    )

    return wrapper
  }

  it('redirects root to /bookmarks', () => {
    const wrapper = setup("/")
    expect(wrapper.find('h2').text()).toBe("Bookmarks")
  })

  it('renders bookmarks index at /bookmarks', () => {
    const wrapper = setup("/bookmarks")
    expect(wrapper.find('h2').text()).toBe("Bookmarks")
  })

  it('renders the bookmark details at /bookmarks/:id', () => {
    const wrapper = setup("/bookmarks/1")
    expect(wrapper.find('h2').text()).toBe("Bookmark 1")
  })

  it('renders tag index at /tags', () => {
    const wrapper = setup("/tags")
    expect(wrapper.find('h2').text()).toBe("Tags")
  })

  it('renders the tag details at /tags/:id', () => {
    const wrapper = setup("/tags/1")
    expect(wrapper.find('h2').text()).toBe("Tag 1")
  })
})
