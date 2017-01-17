require! {
  'async'
  'events' : {EventEmitter}
  'path'
  'prelude-ls' : {filter}
  './service-tester' : ServiceTester
}


class AppTester extends EventEmitter

  ({@app-config, @logger}) ->


  start-testing: ->
    testers = []
    for protection-level of @app-config.services
      for service-role, service-data of @app-config.services[protection-level]
        service-dir = path.join process.cwd!, service-data.location
        testers.push (new ServiceTester {role: service-role, config: {root: service-dir}, @logger})
    async.series [tester.start for tester in testers], (err, exit-codes) ~>
      | err                             =>  @logger.log role: 'exo-test', text: 'Tests failed'; process.exit 1
      | @_contains-non-zero exit-codes  =>  @logger.log role: 'exo-test', text: 'Tests failed'; process.exit 1
      | otherwise                       =>  @logger.log role: 'exo-test', text: 'All tests passed'


  _contains-non-zero: (exit-codes) ->
    exit-codes.filter (> 0)
              .length > 0



module.exports = AppTester
