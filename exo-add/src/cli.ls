require! {
  'chalk' : {cyan, green, red, blue}
  'inquirer'
  'fs'
  'glob'
  'js-yaml' : yaml
  '../../exosphere-shared' : {Logger, templates-path}
  'merge'
  'nitroglycerin' : N
  '../package.json' : {version}
  'path'
  'tmplconv'
  'yaml-cutter'
}

add = ->

  if process.argv[2] is "help"
    help!
    return

  console.log 'We are about to add a new Exosphere service to the application!\n'

  {data, questions} = parse-command-line process.argv
  try
    app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  catch
    if e.code is 'ENOENT'
      console.log red "Error: application.yml not found. Please run this command in the root directory of an Exosphere application."
      process.exit!
    throw e
  inquirer.prompt(questions).then (answers) ->
    data := merge data, answers
    src-path = path.join templates-path, 'add-service' data.template-name
    target-path = path.join process.cwd!, data.service-name
    data.app-name = app-config.name
    tmplconv.render(src-path, target-path, {data}).then ->
      options =
        file: 'application.yml'
        root: 'services.public'
        key: data.service-name
        value: {location: "./#{data.service-name}"}
      yaml-cutter.insert-hash options, N ->
        console.log green "\ndone"


# Returns the names of all known service templates
function service-names
  fs.readdir-sync path.join(templates-path, 'add-service')

function help
  help-message =
    """
    \nUsage: #{cyan 'exo-add'} #{blue '[<entity-name>]'}

    Adds a new service to the current application.
    This command must be called in the root directory of the application.

    options: #{blue '[<service-name>] [<template>] [<model>] [<description>]'}
    """
  console.log help-message



# Returns the data the user provided on the command line,
# and a list of questions that the user has to be asked still.
function parse-command-line command-line-args
  data = {}
  questions = []
  [_, _, entity-name, service-name, author, template-name, model-name, description] = command-line-args

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

  if author
    data.author = author
  else
    questions.push do
      message: 'Author:'
      type: 'input'
      name: 'author'
      filter: (input) -> input.trim!
      validator: (input) -> input.length > 0

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



module.exports = add
