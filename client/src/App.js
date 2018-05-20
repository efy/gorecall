import React, { Component } from 'react'
import { Redirect, Switch, Route } from 'react-router'
import logo from './logo.svg'
import 'spectre.css'
import './App.css'

class App extends Component {
  render() {
    return (
      <div className="app">
        <img src={logo} className="app-logo" alt="logo" />
        <main className="app-main">
          <Switch>
            <Route exact path="/">
              <Redirect to="/bookmarks" />
            </Route>

            <Route exact path="/bookmarks">
              <h2>Bookmarks</h2>
            </Route>

            <Route path="/bookmarks/:id">
              <h2>Bookmark 1</h2>
            </Route>

            <Route exact path="/tags">
              <h2>Tags</h2>
            </Route>

            <Route path="/tags/:id">
              <h2>Tag 1</h2>
            </Route>
          </Switch>
        </main>
      </div>
    );
  }
}

export default App;
