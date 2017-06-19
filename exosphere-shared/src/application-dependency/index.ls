require! {
  './exocom': Exocom
  './nats': Nats
  './generic-dependency': GenericDependency
}

build = (config) ->
  switch config.name
  case 'exocom'
    new Exocom config
  case 'nats'
    new Nats config
  default
    new GenericDependency config
    # throw new Error "Unsupport dependency type: #{config.name}"

module.exports = {build}
