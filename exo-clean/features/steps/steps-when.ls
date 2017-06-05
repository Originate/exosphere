require! {
  'dim-console'
  '../../../exosphere-shared' : {call-args}
  'js-yaml' : yaml
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'fs'
}


module.exports = ->


  @When /^running "([^"]*)" in the terminal$/ timeout: 6_000, (command, done) ->
    if process.platform is 'win32' then command += '.cmd'
    @process = new ObservableProcess(call-args(path.join process.cwd!, 'bin', command),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', -> done!


  @When /^adding a file to the "([^"]*)" service$/ (service-name) ->
    app-config = yaml.safe-load fs.read-file-sync(path.join(@app-dir, 'application.yml'), 'utf8')
    service-config = app-config.services[\public][service-name] or app-config.services[\private][service-name]
    fs.write-file-sync path.join(@app-dir, service-config.location, 'test.txt'), 'test'
