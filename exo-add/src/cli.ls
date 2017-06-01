require! {
  'chalk' : {cyan, green, red, blue}
  'inquirer'
  'fs'
  'glob'
  'js-yaml' : yaml
  '../../exosphere-shared' : {ServiceAdder, Logger, templates-path}
  'merge'
  'nitroglycerin' : N
  'path'
  'prelude-ls' : {flatten, reject}
  'tmplconv'
  'yaml-cutter'
}

module.exports = ->

  if process.argv[2] is "help"
    return help!

  console.log 'We are about to add a new Exosphere service to the application!\n'

  {data, questions} = ServiceAdder.parse-command-line process.argv
  try
    app-config = yaml.safe-load fs.read-file-sync('application.yml', 'utf8')
  catch
    if e.code is 'ENOENT'
      console.log red "Error: application.yml not found. Please run this command in the root directory of an Exosphere application."
      process.exit!
    throw e
  inquirer.prompt(questions).then (answers) ->
    data := merge data, answers
    check-for-service {existing-services: get-existing-services(app-config.services), service-role: data.service-role}
    src-path = path.join templates-path, 'add-service' data.template-name
    target-path = path.join process.cwd!, data.service-type
    data.app-name = app-config.name
    pattern = ['**/*.*', '\.dockerignore', 'Dockerfile']
    tmplconv.render(src-path, target-path, {data, pattern}).then ->
      options =
        file: 'application.yml'
        root: "services.#{data.protection-level}"
        key: data.service-role
        value: {location: "./#{data.service-type}"}
      yaml-cutter.insert-hash options, N ->
        console.log green "\ndone"


function help
  help-message =
    """
    \nUsage: #{cyan 'exo add'} #{blue '[<entity-name>]'}

    Adds a new service to the current application.
    This command must be called in the root directory of the application.

    options: #{blue '--service-role=[<service-role>] --service-type=[<service-type>] --template-name=[<template-name>] --model-name=[<model-name>] --protection-level=[<protection-level>] --description=[<description>]'}
    """
  console.log help-message


function check-for-service {service-role, existing-services}
  if existing-services.includes service-role
    console.log red "Service #{cyan service-role} already exists in this application"
    process.exit 1


function get-existing-services services
  existing-services = []
  for protection-level of services
    if services[protection-level]
      existing-services.push Object.keys that
  flatten existing-services
