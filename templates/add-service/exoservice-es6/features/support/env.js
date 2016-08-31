fs = require('fs-extra')


module.exports = function() {

  this.setDefaultTimeout(1000)
  this.exocom && this.exocom.close()
  this.process && this.process.close()
  this.process && this.process.closePort()

}
