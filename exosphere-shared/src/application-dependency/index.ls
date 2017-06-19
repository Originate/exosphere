require! {
  './exocom': Exocom
  './nats': Nats
}

build = (config) ->
  switch config.name
  case 'exocom'
    new Exocom config
  case 'nats'
    new Nats config
  default
    throw new Error "Unsupport dependency type: #{config.name}"

module.exports = {build}
