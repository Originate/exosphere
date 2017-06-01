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
        if service-data.location
          @services.push do
            role: service-role
            location: service-data.location
    @testers = {}
    for service in @services
      @testers[service.role] = new ServiceTester {service.role, service-location: path.join(process.cwd!, service.location), @logger}

    async.series [tester.start for _, tester of @testers], (err, exit-codes) ~>
      | err                             =>  @logger.log role: 'exo-test', text: 'Tests failed'; process.exit 1
      | @_contains-non-zero exit-codes  =>  @logger.log role: 'exo-test', text: 'Tests failed'; process.exit 1
      | otherwise                       =>  @logger.log role: 'exo-test', text: 'All tests passed'


  _contains-non-zero: (exit-codes) ->
    exit-codes.filter (> 0)
              .length > 0



module.exports = AppTester
