require! {
  'chokidar' : {watch}
  'events' : {EventEmitter}
  '../../exosphere-shared' : {DockerHelper}
  'path'
}


# Watches local services for changes and restarts them
class ServiceWatcher extends EventEmitter

  ({@role, @service-location, @env, @logger}) ->
    @docker-config-location = path.join process.cwd!, 'tmp'


  watch: ~>
    /* Ignores any sub-path including dotfiles.
    '[\/\\]' accounts for both windows and unix systems, the '\.' matches a single '.', and the final '.' matches any character. */
    @watcher = watch @service-location, ignore-initial: yes, ignored: [/.*\/node_modules\/.*/, /(^|[\/\\])\../]
      ..on 'add', (added-path) ~>
        @logger.log {role: 'exo-run', text: "Restarting service '#{@role}' because #{added-path} was created"}
        @_restart!
      ..on 'change', (changed-path) ~>
        @logger.log {role: 'exo-run', text: "Restarting service '#{@role}' because #{changed-path} was changed"}
        @_restart!
      ..on 'unlink', (removed-path) ~>
        @logger.log {role: 'exo-run', text: "Restarting service '#{@role}' because #{removed-path} was deleted"}
        @_restart!


  _restart: ->
    @watcher.close!
    cwd = @docker-config-location 
    DockerHelper.kill-container {service-name: @role, cwd, @write}, (exit-code) ~>
      | exit-code => @emit 'error', "Docker failed to kill container #{@role}"
      @write "Docker container stopped"

      DockerHelper.create-new-container {service-name: @role, cwd, @env, @write}, (exit-code) ~>
        | exit-code => @emit 'error', "Docker image failed to rebuild #{@role}"
        @write "Docker image rebuilt"

        DockerHelper.start-container {service-name: @role, cwd, @env, @write}, (exit-code) ~>
          | exit-code => @emit 'error', "Docker container failed to restart #{@role}"
          @watch!
          @logger.log {role: \exo-run, text: "'#{@role}' restarted successfully"}

  
  write: (text) ~>
    @logger.log {@role, text, trim: yes}


module.exports = ServiceWatcher
