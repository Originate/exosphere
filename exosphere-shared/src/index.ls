require! {
  './logger' : Logger
  './docker-compose' : DockerCompose
  './docker-helper' : DockerHelper
  './call-args'
  './compile-service-routes'
  './normalize-path'
  './kill-child-processes'
  'path'
}


module.exports = {
  call-args
  compile-service-routes
  Logger
  DockerCompose
  DockerHelper
  example-apps-path: path.join(__dirname, '..' 'example-apps')
  normalize-path
  kill-child-processes
  templates-path: path.join(__dirname, '..' 'templates')
}
