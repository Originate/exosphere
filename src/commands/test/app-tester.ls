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
    testers = for service-name in Object.keys @app-config.services
      service-dir = path.join process.cwd!, @app-config.services[service-name].location
      new ServiceTester service-name, root: service-dir
        ..on 'output', (data) ~> @emit 'output', data
        ..on 'service-tests-passed', (name) ~> @emit 'service-tests-passed', name
        ..on 'service-tests-failed', (name) ~> @emit 'service-tests-failed', name
    async.series [tester.start for tester in testers], (err, exit-codes) ~>
      | err                             =>  @emit 'all-tests-failed'
      | @_contains-non-zero exit-codes  =>  @emit 'all-tests-failed'
      | otherwise                       =>  @emit 'all-tests-passed'


  _contains-non-zero: (exit-codes) ->
    exit-codes.filter (> 0)
              .length > 0



module.exports = AppTester
