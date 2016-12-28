require! {
  'async'
  'events' : {EventEmitter}
  'path'
  'prelude-ls' : {filter}
  './service-tester' : ServiceTester
}


class AppTester extends EventEmitter

  (@app-config) ->


  start-testing: ->
    testers = []
    for type of @app-config.services
      for service-name, service-data of @app-config.services[type]
        service-dir = path.join process.cwd!, service-data.location
        testers.push (new ServiceTester service-name, root: service-dir
          ..on 'output', (data) ~> @emit 'output', data
          ..on 'error', (err) ~> @emit 'error', err
          ..on 'service-tests-passed', (name) ~> @emit 'service-tests-passed', name
          ..on 'service-tests-failed', (name) ~> @emit 'service-tests-failed', name
          ..on 'service-tests-skipped', (name) ~> @emit 'service-tests-skipped', name)
    async.series [tester.start for tester in testers], (err, exit-codes) ~>
      | err                             =>  @emit 'all-tests-failed'
      | @_contains-non-zero exit-codes  =>  @emit 'all-tests-failed'
      | otherwise                       =>  @emit 'all-tests-passed'


  _contains-non-zero: (exit-codes) ->
    exit-codes.filter (> 0)
              .length > 0



module.exports = AppTester
