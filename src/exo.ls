require! {
  'chalk' : {dim, red}
  'fs'
  '../package.json' : {version}
}


console.log dim "Exosphere SDK #{version}\n"
command-name = process.argv[2]
command-handler-path = "#{__dirname}/commands/#{command-name}/cli.js"
fs.access command-handler-path, (err) ->
  | err  =>  return unknown-command command-name
  require command-handler-path



# Displays help text when the user provides an unknown command
function unknown-command command
  console.log red "Error: unknown command '#{command-name}'\n"
  console.log 'Usage: exo <command> [options]\n'
  console.log 'Available commands are:'
  for dir in fs.readdir-sync './dist/commands'
    console.log "* #{dir}"
  console.log ''
