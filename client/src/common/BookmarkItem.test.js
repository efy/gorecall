import React from 'react'
import { shallow } from 'enzyme'
import BookmarkItem from './BookmarkItem'

it('renders', () => {
  const wrapper = shallow(<BookmarkItem />)
  expect(wrapper.exists()).toBe(true)
})
