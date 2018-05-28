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
      }
    }

    if (instance.el !== element) {
      instance.el = element
    }

    return instance
  }

  autoLoadVisible () {
    let el = document.getElementById('autoload')
    return (el.offsetParent === null)
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
