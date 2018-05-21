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

import BookmarksPage from './pages/Bookmarks'
import BookmarkPage from './pages/Bookmark'
import NewBookmarkPage from './pages/NewBookmark'
import TagsPage from './pages/Tags'
import TagPage from './pages/Tag'
import NewTagPage from './pages/NewTag'
import SettingsPage from './pages/Settings'
import NotFoundPage from './pages/NotFound'

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
                <BookmarksPage />
              </Route>

              <Route exact path="/bookmarks/search">
                <BookmarksPage />
              </Route>

              <Route exact path="/bookmarks/new">
                <NewBookmarkPage />
              </Route>

              <Route exact path="/bookmarks/:id">
                <BookmarkPage />
              </Route>

              <Route exact path="/tags">
                <TagsPage />
              </Route>

              <Route exact path="/tags/new">
                <NewTagPage />
              </Route>

              <Route exact path="/tags/:id">
                <TagPage />
              </Route>

              <Route path="/settings">
                <SettingsPage />
              </Route>

              <Route>
                <NotFoundPage />
              </Route>
            </Switch>
          </main>
        </div>
      </div>
    );
  }
}

export default App;
