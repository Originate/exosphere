require! {
  'path'
  'require-yaml'
}

module.exports = (app-config, base-path) ->
  service-messages = []
  for type of app-config.services
    for service-name, service-data of app-config.services["#{type}"]
      service-config = require path.join(base-path ? process.cwd!, service-data.location, 'service.yml')
      service-messages.push do
        {
          name: service-name
          receives: service-config.messages.receives
          sends: service-config.messages.sends
        }
  service-messages
