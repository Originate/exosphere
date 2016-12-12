require! {
  'path'
  'fs'
  'js-yaml' : yaml
}

module.exports = (app-config, base-path) ->
  service-messages = []
  for type of app-config.services
    for service-name, service-data of app-config.services["#{type}"]
      service-config = yaml.safe-load fs.read-file-sync(path.join(base-path ? process.cwd!, service-data.location, 'service.yml'), 'utf8')
      service-messages.push do
        {
          name: service-name
          receives: service-config.messages.receives
          sends: service-config.messages.sends
          namespace: service-data.namespace
        }
  service-messages
