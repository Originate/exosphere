const {defineSupportCode} = require('cucumber')

defineSupportCode(function({After, setDefaultTimeout, setWorldConstructor}) {

  setDefaultTimeout(1000)
  After(function (_, done) {
    this.process.kill()
    this.exocom.close(done)
  })

})
