# Converts the given Windows path to a Git Bash path
bash-path = (path) ->
  path.replace /\\/g, '/'
      .replace 'C:', '/c'



module.exports = bash-path
