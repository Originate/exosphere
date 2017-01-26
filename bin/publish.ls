require! {
  'chalk' : {bold, red, cyan, yellow, green}
  'inquirer'
  'jsonfile'
  'semver'
  'shelljs/global'
}

{name, version} = jsonfile.readFileSync 'package.json'

if process.argv.length != 3
  display-help!

level = process.argv[2]
display-help! if level in ['-h, --help']

target-version = semver.inc version, level
unless target-version
  console.log "\n#{bold red 'Error:'} #{bold cyan level} #{red 'is not a valid version increment'}"
  display-help!

console.log "\nYou are about to bump #{green bold name} version #{cyan bold version} up to #{cyan bold target-version}\n"

question =
  type: 'list'
  name: 'continue'
  message: 'Are you sure?'
  choices: ['yes', 'no']
inquirer.prompt([question]).then (answer) ->
  if answer.continue == 'no'
    console.log '\nAborting ...\n'
    process.exit!

  console.log!
  check-exocom-dependencies!
  check-npm-dependencies!
  run-tests!


published-directories =
  'cli'
  'exo-add'
  'exo-clone'
  'exo-create'
  'exo-deploy'
  'exo-lint'
  'exo-run'
  'exo-setup'
  'exo-sync'
  'exo-test'
  'exosphere-shared'


function check-exocom-dependencies
  console.log \exocom


function check-npm-dependencies
  console.log cyan "Checking npm dependencies..."
  run-command-in-subdirs do
    command: './node_modules/.bin/update-check'
    command-message: 'Checking dependencies'
    passing-message: 'dependencies up to date'
    failing-message: 'dependencies not up to date'
    options: silent:true
  console.log '\n'


function run-tests
  console.log cyan "Running tests in all subprojects..."
  run-command-in-subdirs do
    command: './bin/spec'
    command-message: 'Running tests'
    passing-message: 'tests passing'
    failing-message: 'tests failing'
  console.log '\n'


function run-command-in-subdirs {command, command-message, passing-message, failing-message, options}
  for directory in published-directories
    console.log "  #{command-message} in subproject #{cyan directory}"
    cd directory
    if exec(command, options).code > 0
      console.log red "    '#{directory}' #{failing-message}"
      process.exit 1
    console.log green "    '#{directory}' #{passing-message}"
    cd '..'


function display-help
  console.log "\nUsage:\n\n  #{bold 'publish <patch|minor|major>'}\n" #update depending on how publish script is called
  process.exit 1
