require! {
  'abbrev'
  'chalk' : {red}
  'fs'
  'path'
}


entity-name = process.argv[2]
return missing-entity! unless entity-name
command-handler-path = "#{__dirname}/entities/#{abbrev(entity-names!)[entity-name]}/cli.js"
fs.access command-handler-path, (err) ->
  | err  =>  return unknown-command entity-name
  require command-handler-path


function missing-entity
  console.log red "Error: missing entity for 'create' command\n"
  print-usage!


function unknown-command entity
  console.log red "Error: cannot create '#{entity-name}'\n"
  print-usage!


function print-usage
  console.log 'Usage: exo create [<entity>] [<name>] [<template>] [<model>] [<description>]\n'
  console.log 'Available entities are:'
  for entity in entity-names!
    console.log "* #{entity}"
  console.log!


function entity-names
  fs.readdir-sync path.join(__dirname, 'entities')
