require! {
  '../..' : ExoService
  'path'
  'prelude-ls' : {any}
  'wait' : {wait-until}
}


# Provides steps for testing against the JS API
ApiWorld = !->

  @create-exoservice-instance = ({role, exocom-port}, done) ->
    @exoservice = new ExoService {
      root: path.join('features', 'example-services', role)
      exocom-host: "localhost"
      exocom-port
      role
      }
      ..on 'online', (port) -> done!
      ..connect!


  @remove-register-service-message = (@exocom, done) ->
    wait-until (~> @exocom.received-messages.length), 10, ~>
      @exocom.reset! if @exocom.received-messages |> any (.name is 'exocom.register-service')
      done!

module.exports = ApiWorld
