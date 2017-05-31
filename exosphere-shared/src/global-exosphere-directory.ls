require! {
  'path'
  'os'
}

module.exports = (app-name, dependency-name) ->
  path.join os.homedir!, '.exosphere', app-name, dependency-name, 'data'
