import React from 'react';
import { shallow } from 'enzyme';
import renderer from 'react-test-renderer'
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

  it('renders correctly', () => {
    const tree = renderer.create(<App />).toJSON()
    expect(tree).toMatchSnapshot()
  })
})
