defineSupportCode = require('cucumber').defineSupportCode;
fs = require('fs-extra')

defineSupportCode(function({After, setDefaultTimeout, setWorldConstructor}) {

  setDefaultTimeout(1000)
  setWorldConstructor(World)
  After(function () {
  	this.exocom && this.exocom.close()
  	this.process && this.process.close()
  })

});
