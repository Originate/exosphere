require! {
  'fs-extra' : fs
}


module.exports = ->

  @set-default-timeout 1000
  @exocom?.close!
  @process?.close!
  @process?.close-port!
