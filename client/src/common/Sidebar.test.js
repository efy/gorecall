import React from 'react'
import { shallow } from 'enzyme'
import Sidebar from './Sidebar'

it('renders', () => {
  const wrapper = shallow(<Sidebar />)
  expect(wrapper.exists()).toBe(true)
})

