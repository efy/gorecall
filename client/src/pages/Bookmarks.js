import React from 'react'
import BookmarkItem from '../common/BookmarkItem'

const bookmarksMock = [
  {
    id: 1,
    title: 'Example bookmark item',
    url: 'http://example.com/example-bookmark-item',
    domain: 'example.com'
  },
  {
    id: 1,
    title: 'Example bookmark item',
    url: 'http://example.com/example-bookmark-item',
    domain: 'example.com'
  },
  {
    id: 1,
    title: 'Example bookmark item',
    url: 'http://example.com/example-bookmark-item',
    domain: 'example.com'
  },
  {
    id: 1,
    title: 'Example bookmark item',
    url: 'http://example.com/example-bookmark-item',
    domain: 'example.com'
  },
  {
    id: 1,
    title: 'Example bookmark item',
    url: 'http://example.com/example-bookmark-item',
    domain: 'example.com'
  }
]

class Bookmarks extends React.Component {
  render() {
    return (
      <div>
        <h2>Bookmarks</h2>
        {bookmarksMock.map((bookmark, idx) => (
          <BookmarkItem {...bookmark} key={idx} />
        ))}
      </div>
    )
  }
}

export default Bookmarks
