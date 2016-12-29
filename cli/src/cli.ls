require! {
  'abbrev'
  'chalk' : {red}
  '../../exo-add' : add
  '../../exo-clone' : clone
  '../../exo-create' : create
  '../../exo-deploy' : deploy
  '../../exo-lint' : lint
  '../../exo-run' : run
  '../../exo-setup' : setup
  '../../exo-sync' : sync
  '../../exo-test' : test
  'fs'
  'prelude-ls' : {map}
  '../package.json' : pkg
  'path'
  'update-notifier'
}

update-notifier({pkg}).notify!

commands = do
  add: add
  clone: clone
  create: create
  deploy: deploy
  lint: lint
  run: run
  setup: setup
  sync: sync
  test: test

global.templates-path = path.resolve process.cwd!, '../../exosphere-shared/templates'
global.example-apps-path = path.resolve process.cwd!, '../../exosphere-shared/example-apps'

command-name = process.argv[2]
full-command-name = complete-command-name command-name
if command-name is \version
  console.log "Exosphere version #{pkg.version}"
else if not command-name
  missing-command!
else if not full-command-name
  unknown-command command-name
else
  process.argv.shift!
  commands[full-command-name]!


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
  Object.keys commands
