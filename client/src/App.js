import React, { Component } from 'react'
import { Redirect, Switch, Route, Link } from 'react-router-dom'
import logo from './logo.svg'
import 'spectre.css'
import './App.css'

class Sidebar extends Component {
  render() {
    return (
      <header className="app-sidebar">
        <img src={logo} className="app-logo" alt="logo" />
        <Link to="/bookmarks">All</Link>
        <Link to="/bookmarks/search">Search</Link>
        <Link to="/bookmarks/new">Add Link</Link>

        <Link to="/tags">All</Link>
        <Link to="/tags/search">Search</Link>
        <Link to="/tags/new">Add Link</Link>
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
            </Switch>
          </main>
        </div>
      </div>
    );
  }
}

export default App;
