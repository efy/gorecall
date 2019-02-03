import React from 'react'
import { Link } from 'react-router-dom'
import moment from 'moment'
import { Trash2 } from 'react-feather'
import PropTypes from 'prop-types'
import placeholder from '../images/placeholder_favicon.png'
import './BookmarkItem.css'

class BookmarkItem extends React.Component {
  static propTypes = {
    icon: PropTypes.string,
    url: PropTypes.string,
    domain: PropTypes.string,
    created: PropTypes.oneOfType([
      PropTypes.instanceOf(Date),
      PropTypes.string
    ]),
    title: PropTypes.string
  }

  static defaultProps = {
    created: new Date()
  }

  timeFormatted = () => {
    const time = moment(this.props.created)
    return time.fromNow()
  }

  render() {
    return (
      <div className="bookmark-item">
        <div class="rc-bm-favicon">
          <img width="18" height="18" alt={`${this.props.domain} favicon`} src={this.props.icon || placeholder} />
        </div>
        <div>
          <div class="rc-bm-title text-ellipsis">
            <a href={this.props.url} target="_blank" rel="noopener">
              { this.props.title }
            </a>
          </div>
          <div class="rc-bm-details">
            <time>
              { this.timeFormatted() }
            </time>
            {" "} • {" "}
            <Link to={"/bookmarks/" + this.props.id}>
              show
            </Link>
            {" "} • {" "}
            <a href={this.props.domain} rel="noopener" target="_blank">
              { this.props.domain }
            </a>
          </div>
        </div>
        <div class="bookmark-item-actions">
          <button class="btn btn-sm btn-default">
            <Trash2 size={14} />
          </button>
        </div>
      </div>
    )
  }
}

export default BookmarkItem
