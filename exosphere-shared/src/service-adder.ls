require! {
  'fs'
  'minimist'
  'path'
  'prelude-ls': {find}
}

# Performs command line parsing for adding a service in `exo create service` and `exo add`
class ServiceAdder

  @get-child-directory-choices = (dir) ->
    children = fs.readdir-sync dir
    for child in children
      name: child, value: path.join dir, child

  @get-template-path-choices = (cwd) ->
    roles = @get-child-directory-choices path.join(__dirname, '..' 'templates' 'add-service')
    custom-path = path.join(cwd, '.exosphere' 'templates' 'add-service')
    has-custom = false
    try
      fs.access-sync custom-path
      has-custom = true
    if has-custom
      roles = roles.concat @get-child-directory-choices custom-path
    roles

  @parse-command-line = (cwd, command-line-args) ->
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

    template-path-choices = @get-template-path-choices cwd
    if template-name
      data.template-path = template-path-choices |> find (.name is template-name) |> (.value)
    else
      questions.push do
        message: 'Template:'
        type: 'list'
        name: 'templatePath'
        choices: template-path-choices

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
