require! {
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
}


class ServiceSetup extends EventEmitter

  (@name, @config) ->
    @service-config = yaml.safe-load fs.readFileSync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    @emit 'start', @name

    new ObservableProcess(@service-config.setup,
                          cwd: @config.root,
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code is 0 => @emit 'finished', @name
        | otherwise      => @emit 'error', @name, exit-code
        done!

  write: (text) ~>
    @emit 'output', {@name, text}




module.exports = ServiceSetup
