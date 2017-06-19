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

module.exports = {build}
