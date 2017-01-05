require! {
  'async'
  'chalk' : {red, green}
  'child_process'
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'path'
  'rimraf'
  './service-cloner' : ServiceCloner
  'prelude-ls' : {flatten}
}

class AppCloner extends EventEmitter

  ({@repository, @logger}) ->


  start: ->
    if @git-clone-app process.cwd!, @repository.origin then return @logger.log name: 'exo-clone', text: red "Error: cloning #{@repository.name} failed"
    @logger.log name: \exo-clone', text: "#{@repository.name} Application cloned into #{@repository.path}"
    @verify-is-exo-app (err) ~>
      | err  =>  return @logger.log name: 'exo-clone', text: red "Error: application could not be verified.\n #{err}"
      config-path = path.join @repository.path, 'application.yml'
      @app-config = yaml.safe-load fs.read-file-sync(config-path, 'utf8')
      @logger.set-colors Object.keys(@app-config.services)
      cloners = for type of @app-config.services
        for service-name, service-data of @app-config.services[type]
          service-dir = path.join @repository.path, service-data.local
          service-origin = service-data.origin
          new ServiceCloner {name: service-name, config: {root: @repository.path, path: service-dir, origin: service-origin}, @logger}
      async.series [cloner.start for cloner in flatten cloners when cloner.config.origin], (err, exit-codes) ~>
        | err or @_contains-non-zero exit-codes  =>  @logger.log name: 'exo-clone', text: red "Some services failed to clone or were invalid Exosphere services.\nFailed"
        | otherwise                              =>  @logger.log name: 'exo-clone', text: green "All services successfully cloned.\nDone"
        if err or @_contains-non-zero exit-codes then @remove-dir @repository.path


  _log: (text) ~>
    @logger.log name: 'exo-clone', text: text, trim: yes


  _contains-non-zero: (exit-codes) ->
    exit-codes.filter (> 0)
              .length > 0


  # git commands log all output to stderr regardless of exit code.
  git-clone-app: (cwd, origin) ->
    output = child_process.spawn-sync("git", "clone #{origin}".split(' '),
                                     cwd: cwd,
                                     stdio: [1,2])
    switch output.status
    | 0  =>  @_log output.stderr.to-string!.trim!
    | _  =>  @_log red output.stderr.to-string!.trim!
    output.status


  remove-dir: (dir-path) ->
    rimraf dir-path, (error) ->
      | error  =>  @_log "Could not remove #{dir-path}"


  verify-is-exo-app: (callback) ->
    try
      fs.access-sync path.join @repository.path, 'application.yml'
    catch err
      @remove-dir @repository.path
    finally
      callback err



module.exports = AppCloner
