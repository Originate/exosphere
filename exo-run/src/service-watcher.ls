require! {
  'chokidar' : {watch}
  '../../exosphere-shared' : {DockerHelper}
}


# Watches local services for changes and restarts them
class ServiceWatcher

  ({@role, @service-location, @env, @logger}) ->


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
    DockerHelper.kill-container {service-name: @role, @write}, (exit-code) ~>
      | exit-code => @write "Docker failed to kill container #{@role}"; process.exit exit-code
      @write "Docker container stopped"

      DockerHelper.create-new-container {service-name: @role, @env, @write}, (exit-code) ~>
        | exit-code => @write "Docker image failed to rebuild"; process.exit exit-code
        @write "Docker image rebuilt"

        DockerHelper.start-container {service-name: @role, @env, @write}, (exit-code) ~>
          | exit-code => @write "Docker container failed to restart"; process.exit exit-code
          @logger.log {role: \exo-run, text: "'#{@role}' restarted successfully"}

  
  write: (text) ~>
    @logger.log {@role, text, trim: yes}


module.exports = ServiceWatcher
