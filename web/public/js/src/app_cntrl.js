goog.module('controllers.app')

let instance = null

class AppCntrl extends Silica.Controllers.Base {
  constructor (element = document.createElement('div')) {
    if (instance === null) {
      super(element)
      instance = this

    }
    if (instance.el !== element) {
      instance.el = element
    }

    return instance
  }
}

exports = AppCntrl
