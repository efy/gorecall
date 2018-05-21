// setup file
import { configure } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'

configure({ adapter: new Adapter() })

expect.extend({
  toExistIn(finder, wrapper) {
    const pass = wrapper.find(finder).exists();
    const message = () => `Expected ${finder} ${this.isNot ? 'not ' : ''}to exist`
    return { message, pass }
  }
})
