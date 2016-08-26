require! {
  'abbrev'
  'chalk' : {red}
  'fs'
  '../package.json' : pkg
  'path'
  'update-notifier'
}

update-notifier({pkg}).notify!

command-name = process.argv[2]
return missing-command! unless command-name
command-handler-path = "#{__dirname}/commands/#{abbrev(command-names!)[command-name]}/cli.js"
fs.access command-handler-path, (err) ->
  | err  =>  return unknown-command command-name
  require command-handler-path


function missing-command
  console.log red "Error: missing command\n"
  print-usage!


# Displays help text when the user provides an unknown command
function unknown-command command
  console.log red "Error: unknown command '#{command}'\n"
  print-usage!


function print-usage
  console.log 'Usage: exo <command> [options]\n'
  console.log 'Available commands are:'
  for command in command-names!
    if command is 'add'
      console.log "* add [<service-name>] [<template-name>] [<model-name>] [<description>]"
    else
      console.log "* #{command}"
  console.log!


function command-names
  fs.readdir-sync path.join(__dirname, 'commands')
