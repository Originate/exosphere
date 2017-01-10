require! {
  'async'
  'chalk' : {red}
  'events' : {EventEmitter}
  './service-syncer' : ServiceSyncer
  'prelude-ls' : {flatten, filter}
}


# Pulls Git updates for all services of the application
class AppSyncer extends EventEmitter

  ({@app-config, @logger}) ->


  start-syncing: ->
    syncers = for protection-level of @app-config.services
      for service-role, service of @app-config.services[protection-level] ? {}
        new ServiceSyncer {role: service-role, config: {root: service.location}, @logger}


    async.parallel [syncer.start for syncer in flatten syncers], (err, exit-codes) ~>
      | @_count-errors exit-codes  =>  @logger.log name: \exo-sync, text: "Some services failed to sync"
      | otherwise                  =>  @logger.log name: \exo-sync, text: "Sync successful"


  _count-errors: (exit-codes) ->
    exit-codes.filter (code) -> code is not 0
              .length



module.exports = AppSyncer
