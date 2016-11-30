require! {
  'async'
  'chalk' : {red}
  'child_process'
  'events' : {EventEmitter}
  'fs'
  'js-yaml' : yaml
  'path'
  'rimraf'
  './service-cloner' : ServiceCloner
}

class AppCloner extends EventEmitter

  (@repository) ->


  start: ->
    if @git-clone-app process.cwd!, @repository.origin then return @emit 'app-clone-failed'
    @emit 'app-clone-success'
    @verify-is-exo-app (err) ~>
      | err  =>  return @emit 'app-verification-failed', err
      config-path = path.join @repository.path, 'application.yml'
      @app-config = yaml.safe-load fs.read-file-sync(config-path, 'utf8')
      @emit 'app-config-ready', @app-config
      cloners = for service-name in Object.keys @app-config.services
        service-dir = path.join @repository.path, @app-config.services[service-name].local
        service-origin = @app-config.services[service-name].origin
        new ServiceCloner service-name, root: @repository.path, path: service-dir, origin: service-origin
          ..on 'service-clone-success', (name) ~> @emit 'service-clone-success', name
          ..on 'service-clone-fail', (name) ~> @emit 'service-clone-fail', name
          ..on 'service-invalid', (name) ~> @emit 'service-invalid', name
          ..on 'output', (data) ~> @emit 'output', data
      async.series [cloner.start for cloner in cloners when cloner.config.origin], (err, exit-codes) ~>
        | err                             =>  @emit 'service-clones-failed'
        | @_contains-non-zero exit-codes  =>  @emit 'service-clones-failed'
        | otherwise                       =>  @emit 'all-clones-successful'
        if err or @_contains-non-zero exit-codes then @remove-dir @repository.path


  _log: (text) ~>
    @emit 'output', name: 'exo-clone', text: text, trim: yes


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
