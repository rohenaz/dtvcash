goog.module('router')

class Router {
  constructor () {
    console.log('%cHey there!')
    console.log(navigator.userAgent);
    this.hash = null
  }

  route (hash) {
    console.log('routing', hash)
  }
}

exports = Router
