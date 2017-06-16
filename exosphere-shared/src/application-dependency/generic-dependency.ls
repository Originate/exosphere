class GenericDependency

  ({@version}) ->

  get-env-variables: ->
    {}


  get-service-env-variables: ->
    GENERIC_HOST: @_get-container-name!


  get-docker-config: (app-config, done) ->
    done null, do
      "#{@_get-container-name!}":
        image: "generic:#{@version}"
        container_name: @_get-container-name!


  get-online-text: ->
    'Listening for route connections'


  _get-container-name: ->
    "generic#{@version}"

module.exports = GenericDependency
