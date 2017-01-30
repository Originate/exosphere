require! {
  './app-cloner' : AppCloner
  'chalk' : {red, cyan, blue}
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'path'
  'prelude-ls' : {flatten}
}

clone = ->

  if process.argv[2] is "help"
    return help!

  console.log 'We are going to clone an Exosphere application!\n'

  [_, _, repo-origin] = process.argv
  return missing-origin! unless repo-origin
  repository = repo-info repo-origin

  logger = new Logger

  new AppCloner {repository, logger}
    ..on 'done', -> logger.log role: 'exo-clone', text: 'Done'
    ..start!


function repo-info origin
  repo-name = path.basename origin, '.git'
  repo-path = path.join process.cwd!, repo-name
  repo =
    name: repo-name
    origin: origin
    path: repo-path


function missing-origin
  console.log red "Error: missing repository origin"
  print-usage!


function print-usage
  console.log 'Usage: exo clone <origin>\n'


function help
  help-message =
    """
    \nUsage: #{cyan 'exo clone'} #{blue '[<repository>]'}

    Clones an exosphere application hosted on git, including all required services
    This command should be called from an empty working directory.
    """
  console.log help-message

module.exports = clone
