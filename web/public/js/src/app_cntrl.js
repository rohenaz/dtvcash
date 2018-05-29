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