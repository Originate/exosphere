require! {
  './app-cloner' : AppCloner
  'chalk' : {cyan, green, red, yellow}
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger}
  'path'
  'prelude-ls' : {flatten}
}

clone = ->

  console.log 'We are going to clone an Exosphere application!\n'

  [_, _, repo-origin] = process.argv
  return missing-origin! unless repo-origin
  repository = repo-info repo-origin

  logger = new Logger

  new AppCloner {repository, logger}
    ..on 'service-clone-fail', (name) -> logger.log name: name, text: red "Service cloning failed"
    ..on 'service-invalid', (name) -> logger.log name: name, text: red "#{name} is an invalid service"
    ..on 'service-clones-failed', -> logger.log name: 'exo-clone', text: red "Some services failed to clone or were invalid Exosphere services.\nFailed"
    ..on 'all-clones-successful', -> logger.log name: 'exo-clone', text: green "Services successfully cloned.\nDone"
    ..on 'done', -> logger.log name: 'exo-clone', text: 'Done'
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



module.exports = clone
