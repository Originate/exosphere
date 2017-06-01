require! {
  '../compile-service-routes'
}


class ExocomDependency

  ({@version}) ->

  get-env-variables: ->
    EXOCOM_PORT: process.env.EXOCOM_PORT or 80


  get-service-env-variables: ->
    EXOCOM_HOST: @_get-container-name!
    EXOCOM_PORT: '$EXOCOM_PORT'


  get-docker-config: (app-config) ->
    "#{@_get-container-name!}":
      image: "originate/exocom:#{@version}"
      command: 'bin/exocom'
      container_name: @_get-container-name!
      environment:
        ROLE: 'exocom'
        PORT: '$EXOCOM_PORT'
        SERVICE_ROUTES: compile-service-routes app-config |> JSON.stringify |> (.replace /"/g, '')


  _get-container-name: ->
    "exocom#{@version}"

module.exports = ExocomDependency
