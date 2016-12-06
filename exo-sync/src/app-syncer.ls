require! {
  'async'
  'chalk' : {red}
  'events' : {EventEmitter}
  './service-syncer' : ServiceSyncer
  'prelude-ls' : {filter}
}


# Pulls Git updates for all services of the application
class AppSyncer extends EventEmitter

  (@app-config) ->


  start-syncing: ->
    syncers = []
    for type of @app-config.services
      for service-name, service of @app-config.services[type] ? {}
        syncers.push (new ServiceSyncer service-name, root: service.location
          ..on 'output', ~> @emit 'output', &0)


    async.parallel [syncer.start for syncer in syncers], (err, exit-codes) ~>
      | @_count-errors exit-codes  =>  @emit 'sync-failed'
      | otherwise                  =>  @emit 'sync-success'


  _count-errors: (exit-codes) ->
    exit-codes.filter (code) -> code is not 0
              .length



module.exports = AppSyncer
