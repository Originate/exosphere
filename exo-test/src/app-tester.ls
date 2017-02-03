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
    @services = []
    for protection-level of @app-config.services
      for service-role, service-data of @app-config.services[protection-level]
        @services.push do
          {
            role: service-role
            location: service-data.location
          }
    @testers = {}
    for service in @services
      @testers[service.role] = new ServiceTester {service.role, config: {root: path.join(process.cwd!, service.location), app-name: @app-config.name}, @logger}

    async.series [tester.start for _, tester of @testers], (err, exit-codes) ~>
      | err                             =>  @logger.log role: 'exo-test', text: 'Tests failed'; process.exit 1
      | @_contains-non-zero exit-codes  =>  @logger.log role: 'exo-test', text: 'Tests failed'; process.exit 1
      | otherwise                       =>  @logger.log role: 'exo-test', text: 'All tests passed'


  _contains-non-zero: (exit-codes) ->
    exit-codes.filter (> 0)
              .length > 0



module.exports = AppTester
