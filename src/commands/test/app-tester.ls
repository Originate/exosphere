require! {
  'async'
  'events' : {EventEmitter}
  'path'
  './service-tester' : ServiceTester
}


class AppTester extends EventEmitter

  (@app-config) ->


  start-testing: ->
    for service-name in Object.keys @app-config.services
      service-dir = path.join process.cwd!, @app-config.services[service-name].location
      (@testers or= {})[service-name] = new ServiceTester service-name, root: service-dir
        ..on 'output', (data) ~> @emit 'output', data
        ..on 'done', (name) ~> @emit 'service-test-done', name
    async.series [tester.start for _, tester of @testers], (err) ~>
      | err  =>  @emit 'all-tests-failed'
      | _    =>  @emit 'all-tests-passed'



module.exports = AppTester
