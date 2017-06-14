require! {
  'async'
  './docker-helper' : DockerHelper
  'fs'
  'js-yaml' : yaml
  'path'
}

module.exports = ({app-config, base-path}, done) ->
  service-routes = []

  # compile internal service routes
  for , services of app-config.services
    for role, service-data of services
      if service-data.location
        service-config = yaml.safe-load fs.read-file-sync(path.join(base-path ? process.cwd!, service-data.location, 'service.yml'), 'utf8')
        service-routes.push do
          {
            role: role
            receives: service-config.messages.receives
            sends: service-config.messages.sends
            namespace: service-data.namespace
          }

  # compile list of external service config
  external-services = []
  for , services of app-config.services
    for role, service-data of services
      if service-data['docker-image']
        external-services.push {image: service-data['docker-image'], file-name: 'service.yml', role, service-data.namespace}

  # compile service.yml in each external Docker containers
  async.map-series external-services, (-> DockerHelper.cat-file &0, &1), (err, external-service-configs) ->
    | err => return done null

    # compile external service routes
    external-service-index = 0
    for external-service in external-services
      service-config = yaml.safe-load external-service-configs[external-service-index]
      service-routes.push do
        {
          role: external-service.role
          receives: service-config.messages.receives
          sends: service-config.messages.sends
          namespace: external-service.namespace
        }
      external-service-index++

    done null, service-routes
