require! {
  'chalk' : {green}
  '../../../exosphere-shared' : {templates-path}
  'inquirer'
  'merge'
  'path'
  'prelude-ls': {empty}
  'tmplconv'
}

application = ->

  console.log 'We are about to create a new Exosphere application'
  {data, questions} = parse-command-line process.argv

  inquirer.prompt(questions).then (answers) ->
    data := merge data, answers
    src-path = path.join templates-path, 'create-app'
    target-path = data['app-name']
    console.log!
    tmplconv.render(src-path, target-path, {data}).then ->
      console.log green "\ndone"


function parse-command-line command-line-args
  data = {}
  questions = []
  [_, _, _, app-name, app-version, exocom-version, ...app-description] = command-line-args

  if app-name
    data['app-name'] = app-name
  else
    questions.push do
      type: 'input'
      name: 'app-name'
      message: 'Name of the application to create:'
      filter: (input) -> input.trim!
      validate: (input) -> input.length > 0

  if not empty app-description
    data['app-description'] = app-description.join ' '
  else
    questions.push do
      type: 'input'
      name: 'app-description'
      message: 'Description:'

  if app-version
    data['app-version'] = app-version
  else
    questions.push do
      type: 'input'
      name: 'app-version'
      message: 'Initial version:'
      default: '0.0.1'

  if exocom-version
    data['exocom-version'] = exocom-version
  else
    questions.push do
      type: 'input'
      name: 'exocom-version'
      message: 'ExoCom version:'
      default: '0.16.0' # eventually automate how we find latest version

  {data, questions}



module.exports = application
