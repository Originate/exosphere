require! {
  'path'
  'fs'
  'js-yaml' : yaml
}

module.exports = (app-config, base-path) ->
  service-routes = []
  for protection-level of app-config.services
    for role, service-data of app-config.services["#{protection-level}"]
      service-config = yaml.safe-load fs.read-file-sync(path.join(base-path ? process.cwd!, service-data.location, 'service.yml'), 'utf8')
      service-routes.push do
        {
          role: role
          receives: service-config.messages.receives
          sends: service-config.messages.sends
          namespace: service-data.namespace
        }
  service-routes
