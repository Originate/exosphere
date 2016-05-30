require! {
  'chalk' : {cyan, green, red}
  'inquirer'
  'fs'
  'js-yaml' : yaml
  '../../logger' : Logger
  'nitroglycerin' : N
  '../../../package.json' : {version}
  'path'
  'tmplconv'
  'yaml-cutter'
}

console.log 'We are about to add a new Exosphere service to the application!\n'

questions =

  * message: 'Name of the service to create:'
    type: 'input'
    name: 'service-name'
    filter: (input) -> input.trim!
    validate: (input) -> input.length > 0

  * message: 'Description:'
    type: 'input'
    name: 'description'
    filter: (input) -> input.trim!

  * message: 'Type:'
    type: 'list'
    name: 'template-name'
    choices: service-names!

inquirer.prompt(questions).then (answers) ->
  src-path = path.join __dirname, '..' '..' '..' 'templates' 'add-service' answers['template-name']
  target-path = path.join process.cwd!, answers['service-name']
  app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  answers['app-name'] = app-config.name
  console.log answers
  console.log!
  tmplconv.render src-path, target-path, {data: answers}, ->
    options =
      file: 'application.yml'
      root: 'services'
      key: answers['service-name']
      value: {location: "./#{answers['service-name']}"}
    yaml-cutter.insert-hash options, N ->
      console.log green "\ndone"


# Returns the names of all known service templates
function service-names
  fs.readdir-sync path.join(__dirname, '..' '..' '..' 'templates' 'add-service')
