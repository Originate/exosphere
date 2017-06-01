require! {
  '../compile-service-routes'
}


class ExocomDependency

  ({@version}) ->

  get-docker-config: (app-config) ->
    exocom:
      image: "originate/exocom:#{@version}"
      command: 'bin/exocom'
      container_name: 'exocom'
      environment:
        ROLE: 'exocom'
        PORT: '$EXOCOM_PORT'
        SERVICE_ROUTES: compile-service-routes app-config |> JSON.stringify |> (.replace /"/g, '')


module.exports = ExocomDependency
