goog.module('controllers.app')

let instance = null

class AppCntrl extends Silica.Controllers.Base {
  constructor (element = document.createElement('div')) {
    if (instance === null) {
      super(element)
      instance = this

      this.loadMore = false
      this.loadingMore = false
      window.onscroll = () => {

        let visible = this.autoLoadVisible()
        if (visible !== this.loadMore) {
          Silica.apply(() => {
            this.loadMore = true
          })
        }
      }
    }

    if (instance.el !== element) {
      instance.el = element
    }

    return instance
  }

  autoLoadVisible () {
    let elem = document.getElementById('autoload')
    return elem ? this.isElementInViewport(elem) : false
  }

  isElementInViewport (el) {
    var rect = el.getBoundingClientRect()

    return (
      rect.top >= 0 &&
      rect.left >= 0 &&
      rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) && /*or $(window).height() */
      rect.right <= (window.innerWidth || document.documentElement.clientWidth) /*or $(window).width() */
    )
  }

  startWebtorrent (el, id, magnet) {
    console.log('starting webtorrnt with id, magnet')
    var client = []
    client[id] = new WebTorrent()
    var torrentId = magnet + '&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com'
    client[id].add(torrentId, (torrent) => {
    // Torrents can contain many files. Let's use the .mp4 file
    var file = torrent.files.find((file) => {
        // mp4 first
        let mp4 = file.name.endsWith('.mp4')
        // mp3 second
        let mp3 = file.name.endsWith('.mp3')
        return mp4 || mp3
    })

    // Display the file by adding it to the DOM. Supports video, audio, image, etc. files
    if(file) {
        file.appendTo('#message-' + id)
    } else {
        let el = document.getElementById('message-' + id)
        el.innerHTML = 'No mp4 found in this torrent'
    }

    })
  }
}

AppCntrl.watchers = {
  'loadMore': function (newVal, oldVal) {
    if (newVal) {
      this.LoadingMore = true
      console.log('autoload visible!')
      let offset = window.location.search.split('offset=')[1] || 0
      let newOffset = parseInt(offset) + 25
      window.location.href = 'http://dtv.cash/feed?offset=' + newOffset
    }
  }
}

exports = AppCntrl
