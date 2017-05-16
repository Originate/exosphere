require! {
  'child_process'
  'chokidar' : {watch}
  'dashify'
  './docker-runner' : DockerRunner
  'events' : {EventEmitter}
  '../../exosphere-shared' : {call-args, DockerHelper}
  'fs'
  'handlebars' : Handlebars
  'js-yaml' : yaml
  'nitroglycerin' : N
  'os'
  'observable-process' : ObservableProcess
  'path'
  'port-reservation'
  'prelude-ls' : {last}
  'shelljs': {mkdir}
}


class ServiceRunner extends EventEmitter

  ({@role, @config, @logger}) ->
    @service-config = yaml.safe-load @service-configuration-content!


  start: (done) ~>
    @docker-config =
      author: @service-config.author
      image: dashify @service-config.type
      app-name: dashify @config.app-name
      start-command: @service-config.startup.command
      start-text: @service-config.startup['online-text']
      cwd: @config.root
      env:
        EXOCOM_PORT: @config.EXOCOM_PORT
        ROLE: @role
      publish: @service-config.docker?.publish
      dependencies: @compile-service-dependencies!

    @docker-runner = new DockerRunner {@role, @docker-config, @logger}
        ..start-service!
        ..on 'online', -> done?!
        ..on 'error', (message) ~> @emit 'error', error-message: message
        /* Ignores any sub-path including dotfiles.
        '[\/\\]' accounts for both windows and unix systems, the '\.' matches a single '.', and the final '.' matches any character. */
    @watcher = watch @config.root, ignore-initial: yes, ignored: [/.*\/node_modules\/.*/, /(^|[\/\\])\../]
      ..on 'add', (added-path) ~>
        @logger.log {role: 'exo-run', text: "Restarting service '#{@role}' because #{added-path} was created"}
        @restart!
      ..on 'change', (changed-path) ~>
        @logger.log {role: 'exo-run', text: "Restarting service '#{@role}' because #{changed-path} was changed"}
        @restart!
      ..on 'unlink', (removed-path) ~>
        @logger.log {role: 'exo-run', text: "Restarting service '#{@role}' because #{removed-path} was deleted"}
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
          @logger.log {@role, text: "Docker image rebuilt"}
          @start(~> @logger.log role: \exo-run, text: "'#{@role}' restarted successfully")
        | otherwise       =>
          @logger.log {@role, text: "Docker image failed to rebuild"}
          process.exit exit-code


  compile-service-dependencies: ~>
    dependencies = []
    for dependency-name, dependency-config of @service-config.dependencies or {}
      container-name = "#{@config.app-name}-#{dependency-name}"
      if dependency-config?.docker_flags?
        data-path = path.join os.homedir!, '.exosphere', @config.app-name, dependency-name, 'data'
        mkdir '-p', data-path
        volume = Handlebars.compile(that.volume)({"EXO_DATA_PATH": data-path})
        online-text = that.online_text
        port = that.port
      dependencies.push {container-name, dependency-name, version: dependency-config?.version, volume, online-text, port}
    dependencies


  service-configuration-content: ~>
    | @config.image  =>  DockerHelper.get-config(@config.image)
    | otherwise      =>  fs.read-file-sync(path.join @config.root, 'service.yml')


  shutdown-dependencies: ->
    for dependency in @docker-config.dependencies
      DockerHelper.remove-container "#{dependency.container-name}"


  write: (text) ~>
    @logger.log {@role, text, trim: yes}



module.exports = ServiceRunner
