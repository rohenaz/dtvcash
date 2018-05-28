goog.module('Dtv')

const AppCntrl = goog.require('controllers.app')

let Dtv = {
  AppCntrl: AppCntrl
}

// Export context from closure compiler
window['Dtv'] = Dtv

Silica.setContext('Dtv')
Silica.compile(document)

Silica.apply(() => {
  console.log('must call apply')
})
