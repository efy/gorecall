import React, { Component } from 'react'
import { Redirect, Switch, Route, Link } from 'react-router-dom'
import logo from './logo.svg'
import {
  Search,
  List,
  PlusSquare,
  Tag,
  Settings,
  LogOut
} from 'react-feather'
import 'spectre.css'
import './App.css'

class Sidebar extends Component {
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

class App extends Component {
  render() {
    return (
      <div className="app">
        <div className="app-left">
          <Sidebar />
        </div>
        <div className="app-right">
          <main className="app-main">
            <Switch>
              <Route exact path="/">
                <Redirect to="/bookmarks" />
              </Route>

              <Route exact path="/bookmarks">
                <h2>Bookmarks</h2>
              </Route>

              <Route exact path="/bookmarks/search">
                <h2>Bookmarks search</h2>
              </Route>

              <Route exact path="/bookmarks/new">
                <h2>New bookmark</h2>
              </Route>

              <Route exact path="/bookmarks/:id">
                <h2>Bookmark 1</h2>
              </Route>

              <Route exact path="/tags">
                <h2>Tags</h2>
              </Route>

              <Route exact path="/tags/new">
                <h2>New tag</h2>
              </Route>

              <Route exact path="/tags/:id">
                <h2>Tag 1</h2>
              </Route>

              <Route path="/settings">
                <h2>Settings</h2>
              </Route>

              <Route>
                <h2>
                  404 Not found
                </h2>
              </Route>
            </Switch>
          </main>
        </div>
      </div>
    );
  }
}

export default App;
