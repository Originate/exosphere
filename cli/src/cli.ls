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
  '../../package.json' : pkg
  'path'
  'update-notifier'
}

update-notifier({pkg}).notify!

marked.set-options({renderer: new TerminalRenderer!})

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

command-name = process.argv[2]
full-command-name = complete-command-name command-name
if command-name is \version
  console.log "Exosphere version #{pkg.version}"
else if command-name is \help
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
    * add     Add a service to a pre-existing application
    * clone   Clone an exosphere application hosted on git
    * create  Create a new exosphere application or service
    * deploy  Deploy an application to AWS and DockerHub
    * lint    Ensure all service messages are both sent and received
    * run     Run an exosphere application locally
    * setup   Install dependencies and prepare Docker images for an application
    * sync    Update all application services with their git repositories
    * test    Run feature tests for an application or service

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
