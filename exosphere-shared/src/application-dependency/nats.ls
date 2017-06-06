require! {
  '../compile-service-routes'
}


class NatsDependency

  ({@version}) ->

  get-service-env-variables: ->
    NATS_HOST: @_get-container-name!

  get-docker-config: (app-config) ->
    "#{@_get-container-name!}":
      image: "nats:#{@version}"
      container_name: @_get-container-name!

  _get-container-name: ->
    "nats#{@version}"

module.exports = ExocomDependency
