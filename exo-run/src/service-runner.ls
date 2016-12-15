require! {
  'child_process'
  'chokidar' : {watch}
  'dashify'
  './docker-runner' : DockerRunner
  'events' : {EventEmitter}
  '../../exosphere-shared' : {call-args, DockerHelper}
  'fs'
  'js-yaml' : yaml
  'nitroglycerin' : N
  'observable-process' : ObservableProcess
  'path'
  'port-reservation'
  'prelude-ls' : {last}
}


class ServiceRunner extends EventEmitter

  ({@name, @config, @logger}) ->
    @service-config = yaml.safe-load fs.read-file-sync(path.join(@config.root, 'service.yml'), 'utf8')


  start: (done) ~>
    @docker-config =
      author: @service-config.author
      image: dashify @service-config.title
      start-command: @service-config.startup.command
      start-text: @service-config.startup['online-text']
      cwd: @config.root
      env:
        EXOCOM_PORT: @config.EXOCOM_PORT
        SERVICE_NAME: @name
      publish: @service-config.docker.publish if @service-config.docker
      link: @service-config.docker.link if @service-config.docker

    @docker-runner = new DockerRunner {@name, @docker-config, @logger}
        ..start-service!
        ..on 'online', -> done?!
        ..on 'error', (message) ~> @emit 'error', error-message: message

    @watcher = watch @config.root, ignore-initial: yes, ignored: [/.*\/node_modules\/.*/, /(^|[\/\\])\../] # Ignores any sub path including dotfiles e.g. service/.cli/.swp
      ..on 'add', (added-path) ~>
        @logger.log name: 'exo-run', text: "Restarting service '#{@name}' because #{added-path} was created"
        @restart!
      ..on 'change', (changed-path) ~>
        @logger.log name: 'exo-run', text: "Restarting service '#{@name}' because #{changed-path} was changed"
        @restart!
      ..on 'unlink', (removed-path) ~>
        @logger.log name: 'exo-run', text: "Restarting service '#{@name}' because #{removed-path} was deleted"
        @restart!


  restart: ->
    @docker-runner.docker-container.kill!
    @watcher.close!
    new ObservableProcess(call-args(DockerHelper.get-build-command author: @docker-config.author, name: @docker-config.image),
                          cwd: @config.root,
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code, killed) ~>
        | exit-code is 0  =>
          @logger.log name: @name, text: "Docker image rebuilt"
          @start(~> @logger.log name: \exo-run, text: "'#{@name}' restarted successfully")
        | otherwise       =>
          @logger.log name: @name, text: "Docker image failed to rebuild"
          process.exit exit-code


  write: (text) ~>
    @logger.log {@name, text, trim: yes}



module.exports = ServiceRunner
