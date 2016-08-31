require! {
  'chalk' : {green}
  'inquirer'
  'fs'
  'js-yaml' : yaml
  'merge'
  'nitroglycerin' : N
  'path'
  'tmplconv'
  'yaml-cutter'
}

console.log 'We are about to create a new Exosphere service!\n'

{data, questions} = parse-command-line process.argv
inquirer.prompt(questions).then (answers) ->
  data := merge data, answers
  src-path = path.join __dirname, '..' 'templates' 'add-service' data.template-name
  target-path = path.join process.cwd!, '..' data.service-name
  try
    app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  catch error
    throw new Error "Creation of service '#{data.service-name}' has failed."
  data.app-name = app-config.name
  tmplconv.render(src-path, target-path, {data}).then ->
    options =
      file: 'application.yml'
      root: 'services'
      key: data.service-name
      value: {location: "../#{data.service-name}"}
    yaml-cutter.insert-hash options, N ->
      console.log green "\ndone"



function service-names
  fs.readdir-sync path.join(__dirname, '..' 'templates' 'add-service')

function parse-command-line command-line-args
  data = {}
  questions = []
  [_, _, _, service-name, template-name, model-name, description] = command-line-args

  if service-name
    data.service-name = service-name
  else
    questions.push do
      message: 'Name of the service to create'
      type: 'input'
      name: 'serviceName'
      filter: (input) -> input.trim!
      validate: (input) -> input.length > 0

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
      message: 'Name of the data model (leave blank if no model exists):'
      type: 'input'
      name: 'modelName'
      filter: (input) -> input.trim!

  if description
    data.description = description
  else
    questions.push do
      message: 'Description:'
      type: 'input'
      name: 'description'
      filter: (input) -> input.trim!

  {data, questions}
