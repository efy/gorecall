import React from 'react'
import { shallow } from 'enzyme'
import BookmarkItem from './BookmarkItem'

it('renders', () => {
  const wrapper = shallow(<BookmarkItem />)
  expect(wrapper.exists()).toBe(true)
})

it('displays the time relative to today', () => { 
  let time = new Date()
  time.setDate(time.getDate() - 2)

  const wrapper = shallow(<BookmarkItem created={time} />)
  expect(wrapper.find('time').text()).toBe('2 days ago')
})

it('displays a placeholder favicon', () => {
  const wrapper = shallow(<BookmarkItem />)
  expect(wrapper.find('img').props().src).toContain('placeholder_favicon')
})
