# Converts the given Windows path to a Git Bash path
bash-path = (path) ->
  | path.index-of('\\') == -1   =>  return path
  path.replace(/\\/g, '/')
      .replace('C:', '/c')


module.exports = bash-path
