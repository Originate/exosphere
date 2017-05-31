require! {
  'fs'
  'minimist'
  'path'
}

# Performs command line parsing for adding a service in `exo create service` and `exo add`
class ServiceAdder

  @service-roles = ->
    fs.readdir-sync path.join path.join(__dirname, '..' 'templates' 'add-service')

  @parse-command-line = (command-line-args) ->
    parsed-args = minimist command-line-args
    data = {}
    questions = []
    service-role = parsed-args['service-role']
    service-type = parsed-args['service-type']
    author = parsed-args['author']
    template-name = parsed-args['template-name']
    model-name = parsed-args['model-name']
    protection-level = parsed-args['protection-level']
    description = parsed-args['description']

    if service-role
      data.service-role = service-role
    else
      questions.push do
        message: 'Role of the service to create:'
        type: 'input'
        name: 'serviceRole'
        filter: (input) -> input.trim!
        validate: (input) -> input.length > 0

    if service-type
      data.service-type = service-type
    else
      questions.push do
        message: 'Type of the service to create:'
        type: 'input'
        name: 'serviceType'
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
        message: 'Template:'
        type: 'list'
        name: 'templateName'
        choices: @service-roles!

    if model-name
      data.model-name = model-name
    else
      questions.push do
        message: 'Name of the data model:'
        type: 'input'
        name: 'modelName'
        filter: (input) -> input.trim!

    if protection-level
      data.protection-level = protection-level
    else
      questions.push do
        message: 'Protection level:'
        type: 'list'
        name: 'protectionLevel'
        choices: ['public', 'private']

    {data, questions}

module.exports = ServiceAdder
