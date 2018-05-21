import React, { Component } from 'react'
import { Redirect, Switch, Route, Link } from 'react-router-dom'
import 'spectre.css'
import './App.css'

import Sidebar from './common/Sidebar'

import BookmarksPage from './pages/Bookmarks'
import BookmarkPage from './pages/Bookmark'
import NewBookmarkPage from './pages/NewBookmark'
import TagsPage from './pages/Tags'
import TagPage from './pages/Tag'
import NewTagPage from './pages/NewTag'
import SettingsPage from './pages/Settings'
import NotFoundPage from './pages/NotFound'

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
