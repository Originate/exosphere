require! {
  './exocom': Exocom
}

module.exports = (config) ->
  switch config.type
  case 'exocom'
    new Exocom config
  default
    throw new Error "Unsupport dependency type: #{config.type}"
