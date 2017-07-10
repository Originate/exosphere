# Loads an Exoservice directory and makes it available
# as a convenient JS object

require! {
  'defaults'
  'globby'
  'livescript'
  'path'
}


service-loader = (root = '') ->
  server-file-path = _find-server-path root
  handlers = eval('require')(server-file-path) # 'eval' used to hide require satement from Webpack
  if not handlers.before-all? then handlers.before-all = -> it!
  {handlers}


_find-server-path = (root) ->

  # try files in root directory
  files = globby.sync path.join(process.cwd!, root, 'server.*')
  if files.length > 1 then throw _multi-files-error files
  if files.length is 1 then return files[0]

  # try files in a subdirectory
  files = globby.sync [path.join(process.cwd!, root, '**', 'server.*'), path.join('!**', 'node_modules', '**')]
  if files.length > 1 then throw _multi-files-error files
  if files.length is 1 then return files[0]
  throw new Error "Cannot find server file. It must be named 'server.js' or comparably for your language."


_multi-files-error = (files) ->
  new Error "Multiple server files found: #{files.join ', '}"


module.exports = service-loader
