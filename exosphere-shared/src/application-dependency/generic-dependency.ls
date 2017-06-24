require! {
  'fs-extra' : fs
  'handlebars' : Handlebars
  '../global-exosphere-directory'
  'prelude-ls' : {Obj, map}
}

class GenericDependency

  ({@name, @version, @config}) ->


  get-env-variables: ->
    @config.dependency-environment


  get-service-env-variables: ->
    @config.service-environment


  get-docker-config: (app-config, done) ->
    done null, Obj.compact do
      "#{@_get-container-name!}":
        image: "#{@name}:#{@version}"
        container_name: @_get-container-name!
        ports: @config.ports
        volumes: @_render-volumes @config.volumes, app-config.name


  get-online-text: ->
    @config.online-text


  _get-container-name: ->
    "#{@name}#{@version}"


  _render-volumes: (volumes, app-name) ~>
    if volumes
      data-path = global-exosphere-directory app-name, @name
      fs.ensure-dir-sync data-path
      map ((volume) -> Handlebars.compile(volume)({"EXO_DATA_PATH": data-path})), volumes


module.exports = GenericDependency
