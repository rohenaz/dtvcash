goog.module('Dtv')

const AppCntrl = goog.require('controllers.app')
const Router = goog.require('router')

let Dtv = {
  AppCntrl: AppCntrl,
  Router: Router
}

// Export context from closure compiler
window['Dtv'] = Dtv

Silica.setContext('Dtv')
Silica.setRouter(new Router())
Silica.compile(document)

Silica.apply(() => {
  console.log('must call apply')
})
