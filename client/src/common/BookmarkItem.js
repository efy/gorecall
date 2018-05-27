import React from 'react'
import { Link } from 'react-router-dom'
import PropTypes from 'prop-types'

class BookmarkItem extends React.Component {
  static propTypes = {
    icon: PropTypes.string,
    url: PropTypes.string,
    domain: PropTypes.string,
    created: PropTypes.date,
    title: PropTypes.string
  }

  static defaultProps = {
    created: new Date()
  }

  render() {
    return (
      <div className="bookmarkItem">
        <div class="text-center rc-bm-favicon column col-1">
          <img width="20" height="20" src={this.props.icon} />
        </div>
        <div class="column col-9">
          <div class="rc-bm-title text-ellipsis">
            <a href={this.props.url} target="_blank" rel="noopener">
              { this.props.title }
            </a>
          </div>
          <div class="rc-bm-details">
            <time>
              { this.props.created.toString() }
            </time>
            •
            <Link to={"/bookmarks/" + this.props.id}>
              show
            </Link>
            •
            <a href={this.props.domain} rel="noopener" target="_blank">
              { this.props.domain }
            </a>
          </div>
        </div>
        <div class="column col-2 text-right">
          <button class="btn btn-sm btn-default">Delete</button>
        </div>
      </div>
    )
  }
}

export default BookmarkItem
