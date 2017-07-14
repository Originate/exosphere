require! {
  'abbrev'
  'chalk' : {red}
  'cross-spawn': spawn
  'fs'
  'marked'
  'marked-terminal': TerminalRenderer
  'prelude-ls' : {map}
  '../../package.json' : pkg
  'path'
  'update-notifier'
}

update-notifier({pkg}).notify!
marked.set-options renderer: new TerminalRenderer!

commands = do
  add: "../../exo-add"
  clone: "../../exo-clone"
  clean: "../../exo-clean"
  create: "../../exo-create"
  deploy: "../../exo-deploy"
  lint: "../../exo-lint"
  run: "../../exo-run"
  setup: "../../exo-setup"
  sync: "../../exo-sync"
  "template": "../../exo-template"
  test: "../../exo-test"

go-commands = ['add', 'clean', 'create', 'run', 'template']

command-name = process.argv[2]

if command-name is \version
  return console.log "Exosphere version #{pkg.version}"
if command-name is \help
  process.argv.shift!
  return print-usage! unless process.argv[2]
  process.argv.push "help"

command-name = process.argv[2]
full-command-name = complete-command-name command-name

if not command-name
  missing-command!
else if not full-command-name
  unknown-command command-name
else if full-command-name in go-commands
  args = [full-command-name].concat process.argv.slice(3)
  {error} = spawn.sync get-go-binary-path!, args, stdio: 'inherit'
  throw error if error
else
  process.argv.shift!
  (require commands[full-command-name])!


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
    * add             Add a service to an existing application
    * clean           Remove dangling Docker images and volumes
    * clone           Download the source code of an application
    * create          Create a new application or stand-alone service
    * deploy          Deploy an application to the cloud
    * lint            Verify the correctness of an application
    * run             Run an application locally
    * setup           Prepare a freshly cloned application for running it
    * sync            Download updates for an application from its Git repository
    * template        Manage remote service templates
    * test            Run the tests for an application or service

  Use "exo <command> help" or "exo help <command>" for more information about a specific command.
  """
  console.log marked usage-text


function command-names
  Object.keys commands

function get-go-binary-os
  switch process.platform
    | 'darwin'  => 'darwin'
    | 'linux'   => 'linux'
    | 'win32'   => 'windows'
    | otherwise => throw new Error('Unsupported operating system. Please open an issue with your operating system and system architecture.')

function get-go-binary-architecture
  switch process.arch
    | 'arm'     => 'arm'
    | 'i32'     => '386'
    | 'x64'     => 'amd64'
    | otherwise => throw new Error('Unsupported system architecture system. Please open an issue with your operating system and system architecture.')

function get-go-binary-path
  path.join __dirname, '..', '..', 'go-binaries', "exo-#{get-go-binary-os!}-#{get-go-binary-architecture!}"
