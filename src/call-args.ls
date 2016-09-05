# Returns the arguments for ObservableProcess to call the given command
call-args = (command) ->
  | process.platform is 'win32'  =>  ['cmd' '/c' command]
  | otherwise                    =>  ['bash', '-c' command]


module.exports = call-args
