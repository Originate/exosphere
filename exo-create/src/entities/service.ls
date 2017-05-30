require! {
  'chalk' : {green}
  '../../../exosphere-shared' : {templates-path, ServiceAdder}
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

  {data, questions} = ServiceAdder.parse-command-line process.argv
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



module.exports = service
