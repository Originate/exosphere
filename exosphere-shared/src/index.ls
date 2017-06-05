require! {
  './application-dependency' : ApplicationDependency
  './logger' : Logger
  './docker-compose' : DockerCompose
  './docker-helper' : DockerHelper
  './service-adder' : ServiceAdder
  './call-args'
  './compile-service-routes'
  './global-exosphere-directory'
  './normalize-path'
  './kill-child-processes'
  'path'
}


module.exports = {
  ApplicationDependency
  call-args
  compile-service-routes
  Logger
  DockerCompose
  DockerHelper
  ServiceAdder
  example-apps-path: path.join(__dirname, '..' 'example-apps')
  global-exosphere-directory
  normalize-path
  kill-child-processes
  templates-path: path.join(__dirname, '..' 'templates')
}
