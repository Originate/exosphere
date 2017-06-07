require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'js-yaml' : yaml
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'fs'
  'wait' : {wait}
}


module.exports = ->

  @When /^running "([^"]*)" in the terminal$/ timeout: 6_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!
