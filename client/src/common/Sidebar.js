import React from 'react'
import { Link } from 'react-router-dom'
import logo from '../logo.svg'
import {
  Search,
  List,
  PlusSquare,
  Tag,
  Settings,
  LogOut
} from 'react-feather'

class Sidebar extends React.Component {
  render() {
    return (
      <header className="app-sidebar">
        <img src={logo} className="app-logo" alt="logo" />
        <ul className="nav">
          <li className="divider" data-content="Links"></li>
          <li className="nav-item">
            <Link to="/bookmarks">All</Link>
            <List size={16} />
          </li>
          <li className="nav-item">
            <Link to="/bookmarks/search">Search</Link>
            <Search size={16} />
          </li>
          <li className="nav-item">
            <Link to="/bookmarks/new">Add Link</Link>
            <PlusSquare size={16} />
          </li>

          <li className="divider" data-content="Tags"></li>
          <li className="nav-item">
            <Link to="/tags">All</Link>
            <Tag size={16} />
          </li>
          <li className="nav-item">
            <Link to="/tags/new">Add Tag</Link>
            <PlusSquare size={16} />
          </li>
          <li className="divider"></li>
        </ul>

        <ul className="nav">
          <li className="nav-item">
            <Link to="/settings">Settings</Link>
            <Settings size={16} />
          </li>
          <li className="nav-item">
            <Link to="/logout">Logout</Link>
            <LogOut size={16} />
          </li>
        </ul>
      </header>
    )
  }
}

export default Sidebar
