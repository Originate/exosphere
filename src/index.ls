require! {
  './logger' : Logger
  './call-args'
  './compile-service-messages'
  './normalize-path'
  './kill-child-processes'
  'path'
}


module.exports = {
  call-args
  compile-service-messages
  Logger
  example-apps-path: path.join(__dirname, '..' 'example-apps')
  normalize-path
  kill-child-processes
  templates-path: path.join(__dirname, '..' 'templates')
}
