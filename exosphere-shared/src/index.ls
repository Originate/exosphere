require! {
  './logger' : Logger
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
  call-args
  compile-service-routes
  Logger
  DockerHelper
  ServiceAdder
  example-apps-path: path.join(__dirname, '..' 'example-apps')
  global-exosphere-directory
  normalize-path
  kill-child-processes
  templates-path: path.join(__dirname, '..' 'templates')
}
