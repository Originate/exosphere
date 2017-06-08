class NatsDependency

  ({@version}) ->

  get-env-variables: ->
    {}


  get-service-env-variables: ->
    NATS_HOST: @_get-container-name!


  get-docker-config: (app-config, done) ->
    done null, do
      "#{@_get-container-name!}":
        image: "nats:#{@version}"
        container_name: @_get-container-name!


  get-online-text: ->
    'Listening for route connections'


  _get-container-name: ->
    "nats#{@version}"

module.exports = NatsDependency
