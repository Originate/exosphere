require! {
  './logger' : Logger
  './call-args'
  './normalize-path'
  './kill-child-processes'
  'path'
}


module.exports = {
  call-args
  Logger
  example-apps-path: path.join(__dirname, '..' 'example-apps')
  normalize-path
  kill-child-processes
  templates-path: path.join(__dirname, '..' 'templates')
}
