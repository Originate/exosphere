# Returns the arguments for ObservableProcess to call the given command
call-args = (command) ->
  | process.platform is 'win32'  =>  ['cmd' '/c' command.replace(/\//g, '\\')]
  | otherwise                    =>  ['bash', '-c' command.replace(/\\/g, '/')]


module.exports = call-args
