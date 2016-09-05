require! {
  'abbrev'
  'chalk' : {red}
  'fs'
  'glob'
  'prelude-ls' : {map}
  '../package.json' : pkg
  'path'
  'update-notifier'
}

update-notifier({pkg}).notify!

command-name = process.argv[2]
return missing-command! unless command-name
full-command-name = complete-command-name command-name
return unknown-command command-name unless full-command-name
process.argv.shift!
require path.join __dirname, '..' 'node_modules', "exo-#{full-command-name}", 'dist', 'cli.js'


function complete-command-name command-name
  abbrev(command-names!)[command-name]


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
  glob.sync "#{__dirname}/../node_modules/exo-*"
  |> map (.split '-')
  |> map (parts) -> parts[*-1]
