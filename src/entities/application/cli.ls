require! {
  'chalk' : {green}
  'inquirer'
  'path'
  'tmplconv'
}

console.log 'We are about to create a new Exosphere application!\n'

questions =

  * type: 'input'
    name: 'app-name'
    message: 'Name of the application to create:'
    filter: (input) -> input.trim!
    validate: (input) -> input.length > 0

  * type: 'input'
    name: 'app-description'
    message: 'Description:'

  * type: 'input'
    name: 'app-version'
    message: 'Initial version:'
    default: '0.0.1'

inquirer.prompt(questions).then (answers) ->
  src-path = path.join __dirname, '..' '..' '..' 'templates' 'create-app'
  target-path = answers['app-name']
  console.log!
  tmplconv.render(src-path, target-path, {data: answers}).then ->
    console.log green "\ndone"
