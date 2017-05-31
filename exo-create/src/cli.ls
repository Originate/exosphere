require! {
  'abbrev'
  './entities/application' : application
  'chalk' : {red, blue, cyan}
  'fs'
  'path'
  'prelude-ls' : {map}
  './entities/service' : service
}

entities = do
  application: application
  service: service

module.exports = ->
  if process.argv[2] is \help
    return help!
  entity-name = process.argv[2]
  return missing-entity! unless entity-name
  full-entity-name = abbrev(entity-names!)[entity-name]
  return unknown-command entity-name unless full-entity-name in entity-names!
  entities[full-entity-name]!


function missing-entity
  console.log red "Error: missing entity for 'create' command\n"
  help!


function unknown-command entity
  console.log red "Error: cannot create '#{entity}'\n"
  help!


function help
  help-text = """
  Usage: #{cyan "exo create"} #{blue "[<entity>]"}

  Available entities are:

    * application   Create a new Exosphere application
      - options: #{blue "[<app-name>] [<app-version>] [<exocom-version>] [<app-description>]"}

    * service       Create a new service for this application located in the parent directory
      - options: #{blue "--service-role=[<service-role>] --service-type=[<service-type>] --author=[<author>] --template-name=[<template-name>] --model-name=[<model-name>] --protection-level=[<protection-level>] --description=[<description>]"}
  """
  console.log help-text


function entity-names
  Object.keys entities
