require! {
  'abbrev'
  './entities/application' : application
  'chalk' : {red}
  'fs'
  'path'
  'prelude-ls' : {map}
  './entities/service' : service
}

entities = do
  application: application
  service: service

module.exports = ->
  entity-name = process.argv[2]
  return missing-entity! unless entity-name
  full-entity-name = abbrev(entity-names!)[entity-name]
  return unknown-command entity-name unless full-entity-name in entity-names!
  entities[full-entity-name]!


function missing-entity
  console.log red "Error: missing entity for 'create' command\n"
  print-usage!


function unknown-command entity
  console.log red "Error: cannot create '#{entity}'\n"
  print-usage!


function print-usage
  console.log 'Usage: exo create [<entity>] [<name>] [<template>] [<model>] [<description>]\n'
  console.log 'Available entities are:'
  for entity in entity-names!
    console.log "* #{entity}"
  console.log!


function entity-names
  Object.keys entities
