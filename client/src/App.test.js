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

  it('renders Sidebar', () => {
    const wrapper = shallow(<App />)
    expect(wrapper.find('Sidebar').exists()).toBe(true)
  })
})

describe('Sidebar', () => {

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
    expect('Bookmarks').toExistIn(wrapper)
  })

  it('renders bookmarks index at /bookmarks', () => {
    const wrapper = setup("/")
    expect('Bookmarks').toExistIn(wrapper)
  })

  it('renders bookmarks search at /bookmarks/search', () => {
    const wrapper = setup("/bookmarks/search")
    expect('Bookmarks').toExistIn(wrapper)
  })

  it('renders new bookmark at /bookmarks/new', () => {
    const wrapper = setup("/bookmarks/new")
    expect('NewBookmark').toExistIn(wrapper)
  })

  it('renders the bookmark details at /bookmarks/:id', () => {
    const wrapper = setup("/bookmarks/1")
    expect('Bookmark').toExistIn(wrapper)
  })

  it('renders tag index at /tags', () => {
    const wrapper = setup("/tags")
    expect('Tags').toExistIn(wrapper)
  })

  it('renders new tag at /tags', () => {
    const wrapper = setup("/tags/new")
    expect('NewTag').toExistIn(wrapper)
  })

  it('renders the tag details at /tags/:id', () => {
    const wrapper = setup("/tags/1")
    expect('Tag').toExistIn(wrapper)
  })

  it('renders settings page at /settings/*', () => {
    let wrapper = setup("/settings")
    expect('Settings').toExistIn(wrapper)

    wrapper = setup("/settings/account")
    expect('Settings').toExistIn(wrapper)
  })

  it('renders page not found if no route matches', () => {
    let wrapper = setup("/nomatch")
    expect('NotFound').toExistIn(wrapper)
  })
})
