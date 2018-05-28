goog.module('controllers.app')

let instance = null

class AppCntrl extends Silica.Controllers.Base {
  constructor (element = document.createElement('div')) {
    if (instance === null) {
      super(element)
      instance = this

      this.loadMore = false
      window.onscroll = () => {
        this.loadMore = this.autoLoadVisible()
        console.info('visible?', this.loadMore)
      }
    }

    if (instance.el !== element) {
      instance.el = element
    }

    return instance
  }

  autoLoadVisible () {
    let elem = document.getElementById('autoload')
    return !!(elem.offsetWidth || elem.offsetHeight || elem.getClientRects().length )
  }
}

AppCntrl.watchers = {
  'loadMore': function (newVal, oldVal) {
    if (newVal) {
      console.log('autoload visible!')
    }
  }
}

exports = AppCntrl
