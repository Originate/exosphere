var formatOptions = JSON.stringify({snippetSyntax: 'node_modules/cucumber-snippets-livescript'})

var common = [
  '--compiler', 'ls:livescript',
  '--format', 'progress-bar',
  '--format-options', '\'' + formatOptions + '\'',
  '--require', 'features'
].join(' ')

module.exports = {
  "default": common
}
