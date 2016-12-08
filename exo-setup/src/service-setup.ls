require! {
  'chalk' : {red}
  'events' : {EventEmitter}
  '../../exosphere-shared' : {call-args, normalize-path}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
}


class ServiceSetup extends EventEmitter

  ({@name, @logger, @config}) ->
    @service-config = yaml.safe-load fs.read-file-sync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    @logger.log name: @name, text: "starting setup"
    new ObservableProcess(call-args(normalize-path @service-config.setup),
                          cwd: @config.root,
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code is 0  =>  @logger.log name: @name, text: 'setup finished'
        | otherwise       =>
          @logger.log name: @name, text: "setup failed with exit code #{exit-code}"
          process.exit exit-code
        done!


  write: (text) ~>
    @emit 'output', {@name, text}



module.exports = ServiceSetup
