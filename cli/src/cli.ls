require! {
  '../../package.json' : pkg
}

command-name = process.argv[2]
if command-name is \version # Version handled before all others for performance
  return console.log "Exosphere version #{pkg.version}"

require! {
  'abbrev'
  'chalk' : {red}
  '../../exo-add' : add
  '../../exo-clone' : clone
  '../../exo-create' : create
  '../../exo-deploy' : deploy
  '../../exo-lint' : lint
  '../../exo-run' : run
  '../../exo-setup' : setup
  '../../exo-sync' : sync
  '../../exo-test' : test
  'fs'
  'marked'
  'marked-terminal': TerminalRenderer
  'prelude-ls' : {map}
  'path'
  'update-notifier'
}

update-notifier({pkg}).notify!
marked.set-options renderer: new TerminalRenderer!

commands = do
  add: add
  clone: clone
  create: create
  deploy: deploy
  lint: lint
  run: run
  setup: setup
  sync: sync
  test: test

full-command-name = complete-command-name command-name
if command-name is \help
  process.argv.shift!
  help process.argv[2]
else if not command-name
  missing-command!
else if not full-command-name
  unknown-command command-name
else
  process.argv.shift!
  commands[full-command-name]!


function complete-command-name command-name
  abbrev(command-names!)[command-name]


function missing-command
  console.log red "Error: missing command\n"
  print-usage!


# Displays help text when the user provides an unknown command
function unknown-command command
  console.log red "Error: unknown command '#{command}'\n"
  print-usage!


function print-usage
  usage-text = """
  **Usage: exo <command> [options]**

  Available commands are:
    * add     Add a service to an existing application
    * clone   Download the source code of an application
    * create  Create a new application or stand-alone service
    * deploy  Deploy an application to the cloud
    * lint    Verify the correctness of an application
    * run     Run an application locally
    * setup   Prepare a freshly cloned application for running it
    * sync    Download updates for an application from its Git repository
    * test    Run the tests for an application or service

  Use "exo <command> help" or "exo help <command>" for more information about a specific command.
  """
  console.log marked usage-text


function help command
  return missing-command! unless command
  process.argv.push "help"
  process.argv.shift!
  commands[command]!


function command-names
  Object.keys commands
