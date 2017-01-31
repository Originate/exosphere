require! {
  'chalk' : {green}
  '../../../exosphere-shared' : {templates-path}
  'inquirer'
  'fs'
  'js-yaml' : yaml
  'merge'
  'nitroglycerin' : N
  'path'
  'prelude-ls' : {empty}
  'tmplconv'
  'yaml-cutter'
}

service = ->

  console.log 'We are about to create a new Exosphere service!\n'

  {data, questions} = parse-command-line process.argv
  inquirer.prompt(questions).then (answers) ->
    data := merge data, answers
    src-path = path.join templates-path, 'add-service', data.template-name
    target-path = path.join process.cwd!, '..' data.service-type
    try
      app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
    catch error
      throw new Error "Creation of service '#{data.service-role}' has failed."
    data.app-name = app-config.name
    tmplconv.render(src-path, target-path, {data}).then ->
      options =
        file: 'application.yml'
        root: 'services.public'
        key: data.service-role
        value: {location: "../#{data.service-type}"}
      yaml-cutter.insert-hash options, N ->
        console.log green "\ndone"



function service-roles
  fs.readdir-sync path.join templates-path, 'add-service'


function parse-command-line command-line-args
  data = {}
  questions = []
  [_, _, _, service-role, service-type, author, template-name, model-name, ...description] = command-line-args

  if service-role
    data.service-role = service-role
  else
    questions.push do
      message: 'Role of the service to create'
      type: 'input'
      name: 'serviceRole'
      filter: (input) -> input.trim!
      validate: (input) -> input.length > 0

  if service-type
    data.service-type = service-type
  else
    questions.push do
      message: 'Type of the service to create'
      type: 'input'
      name: 'serviceType'
      filter: (input) -> input.trim!
      validate: (input) -> input.length > 0

  if template-name
    data.template-name = template-name
  else
    questions.push do
      message: 'Template:'
      type: 'list'
      name: 'templateName'
      choices: service-roles!

  if model-name
    data.model-name = model-name
  else
    questions.push do
      message: 'Name of the data model (leave blank if no model exists):'
      type: 'input'
      name: 'modelName'
      filter: (input) -> input.trim!

  if not empty description
    data.description = description.join ' '
  else
    questions.push do
      message: 'Description:'
      type: 'input'
      name: 'description'
      filter: (input) -> input.trim!

  if author
    data.author = author
  else
    questions.push do
      message: 'Author:'
      type: 'input'
      name: 'author'
      filter: (input) -> input.trim!
      validator: (input) -> input.length > 0

  {data, questions}



module.exports = service
