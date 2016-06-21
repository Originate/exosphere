require! {
  'async'
  'events' : {EventEmitter}
  'path'
  './service-tester' : ServiceTester
}


class AppTester extends EventEmitter

  (@app-config) ->


  start-testing: ->
    service-names = Object.keys @app-config.services
    @testers = {}
    for service-name in service-names
      service-dir = path.join process.cwd!, @app-config.services[service-name].location
      @testers[service-name] = new ServiceTester service-name, root: service-dir
        ..on 'output', (data) ~> @emit 'output', data
        ..on 'done', (name) ~> @emit 'service-test-done', name
    async.series [tester.start for _, tester of @testers], (err) ~>
      if err
        @emit 'all-tests-failed'
      else
        @emit 'all-tests-passed'



module.exports = AppTester
