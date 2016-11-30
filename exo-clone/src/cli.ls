require! {
  './app-cloner' : AppCloner
  'chalk' : {cyan, green, red, yellow}
  'js-yaml' : yaml
  'exosphere-shared' : {Logger}
  'path'
}

console.log 'We are going to clone an Exosphere application!\n'

[_, _, repo-origin] = process.argv
return missing-origin! unless repo-origin
repo = repo-info repo-origin

logger = new Logger

new AppCloner repo
  ..on 'output', (data) -> logger.log data
  ..on 'app-config-ready', (app-config) -> logger.set-colors Object.keys(app-config.services)
  ..on 'app-verification-failed', (err) -> logger.log name: 'exo-clone', text: red "Error: application could not be verified.\n" + red err
  ..on 'app-clone-success' -> logger.log name: 'exo-clone', text: "#{repo.name} Application cloned into #{repo.path}"
  ..on 'app-clone-failed', -> logger.log name: 'exo-clone', text: red "Error: cloning #{repo.name} failed"
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
