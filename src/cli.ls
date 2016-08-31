require! {
  'chalk' : {cyan, green, red}
  'inquirer'
  'fs'
  'js-yaml' : yaml
  'exosphere-shared' : {Logger}
  'merge'
  'nitroglycerin' : N
  '../package.json' : {version}
  'path'
  'tmplconv'
  'yaml-cutter'
}

console.log 'We are about to add a new Exosphere service to the application!\n'

{data, questions} = parse-command-line process.argv
inquirer.prompt(questions).then (answers) ->
  data := merge data, answers
  src-path = path.join __dirname, '..' 'node_modules' 'exosphere-shared' 'templates' 'add-service' data.template-name
  target-path = path.join process.cwd!, data.service-name
  try
    app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  catch
    console.log e
    throw e
  data.app-name = app-config.name
  tmplconv.render(src-path, target-path, {data}).then ->
    options =
      file: 'application.yml'
      root: 'services'
      key: data.service-name
      value: {location: "./#{data.service-name}"}
    yaml-cutter.insert-hash options, N ->
      console.log green "\ndone"


# Returns the names of all known service templates
function service-names
  fs.readdir-sync path.join(__dirname, '..' 'node_modules' 'exosphere-shared' 'templates' 'add-service')


# Returns the data the user provided on the command line,
# and a list of questions that the user has to be asked still.
function parse-command-line command-line-args
  data = {}
  questions = []
  [_, _, service-type, service-name, template-name, model-name, description] = command-line-args

  if service-name
    data.service-name = service-name
  else
    questions.push do
      message: 'Name of the service to create:'
      type: 'input'
      name: 'serviceName'
      filter: (input) -> input.trim!
      validate: (input) -> input.length > 0

  if description
    data.description = description
  else
    questions.push do
      message: 'Description:'
      type: 'input'
      name: 'description'
      filter: (input) -> input.trim!

  if template-name
    data.template-name = template-name
  else
    questions.push do
      message: 'Type:'
      type: 'list'
      name: 'templateName'
      choices: service-names!

  if model-name
    data.model-name = model-name
  else
    questions.push do
      message: 'Name of the data model:'
      type: 'input'
      name: 'modelName'
      filter: (input) -> input.trim!

  {data, questions}
