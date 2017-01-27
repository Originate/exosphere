require! {
  'chalk' : {bold, red, cyan, yellow, green, dim}
  'inquirer'
  'jsonfile'
  'replace'
  'semver'
  'shelljs/global'
}

SUBPROJECTS_TO_PUBLISH =
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

{name, version} = jsonfile.readFileSync './package.json'

if process.argv.length != 3
  display-help!

level = process.argv[2]
display-help! if level in ['-h, --help']

target-version = get-target-version version, level

confirm-target-version ->
  # ensure-no-open-files!
  # ensure-on-master!

  check-npm-dependencies!
  update-exocom-dependencies!
  build-subprojects!
  # run-tests!
  # bump-version-number!
  # push-version-number!
  # publish-to-npm!
  # push-exocom-docker-image!
  # push-git-tags!


function check-npm-dependencies
  console.log green "Checking npm dependencies...\n"
  run-command-in-subdirs do
    command: './node_modules/.bin/update-check'
    command-message: 'Checking dependencies'
  console.log!


function update-exocom-dependencies
  update-exorelay-version!


function update-exorelay-version
  replace do
    regex: /"exorelay": "\d+.\d+.\d+"/g
    replacement: "\"exorelay\": \"#{target-version}\""
    paths: ['.']
    recursive: true
    # silent: true


function build-subprojects
  console.log green "Building subprojects...\n"
  run-command-in-subdirs do
    command: './node_modules/.bin/build'
    command-message: 'Building'
  console.log!


function run-tests
  console.log green "Running tests in subprojects...\n"
  run-command-in-subdirs do
    command: './bin/spec'
    command-message: 'Running tests in'
  console.log!


function bump-version-number
  console.log green "Bumping version number...\n"
  run-command "npm version #{level}"
  console.log!


function push-version-number
  console.log green "Pushing version number...\n"
  run-command "git add -u && git commit -m #{target-version} && git push"
  console.log!


function publish-to-npm
  console.log green "Publishing to npm...\n"
  run-command-in-subdirs do
    command: 'npm publish'
    command-message: 'Publishing'
  console.log!


function push-exo-deploy-docker-image #TODO: make sure version number is correct
  console.log green "Pushing exo-deploy image to DockerHub...\n"
  cd 'exo-deploy'
  run-command "docker build --no-cache -t originate/exo-deploy:#{target-version}"
  run-command "docker push originate/exo-deploy:#{target-version}"
  cd '..'


function push-git-tags
  console.log green "Pushing git release tag...\n"
  run-command 'git tag -a v#{target-version} && git push --tags'
  console.log!


function run-command command
  if exec(command).code > 0 then process.exit 1


function run-command-in-subdirs {command, command-message}
  for directory in SUBPROJECTS_TO_PUBLISH
    console.log "#{command-message} subproject #{cyan directory}"
    cd directory
    run-command command
    cd '..'
    console.log!


function get-target-version version, level
  target-version = semver.inc version, level
  unless target-version
    console.log "\n#{bold red 'Error:'} #{bold cyan level} #{red 'is not a valid version increment'}"
    display-help!
  target-version


function confirm-target-version done
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
    done!


function ensure-no-open-files
  open-files = exec "git status --porcelain", silent: yes
  if open-files then console.log red 'Please commit all files before releasing' ; process.exit 1


function ensure-on-master
  current-branch = exec "git rev-parse --abbrev-ref HEAD", silent: yes
  if current-branch.trim! isnt 'master' then console.log red 'You must be on the master branch to publish' ; process.exit 1


function display-help
  console.log "\nUsage:\n\n  #{bold 'publish <patch|minor|major>'}\n"
  process.exit 1
