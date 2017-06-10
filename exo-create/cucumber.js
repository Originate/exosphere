var formatOptions = {
  snippetSyntax: 'node_modules/cucumber-snippets-livescript'
};

var common = [
  '--compiler ls:livescript',
  '-r features',
  '--fail-fast',
  "--format-options '" + JSON.stringify(formatOptions) + "'",
].join(' ');

module.exports = {
  "default": common
};
