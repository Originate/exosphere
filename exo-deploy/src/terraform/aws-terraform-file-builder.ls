require! {
  'async'
  '../../../exosphere-shared' : {compile-service-messages}
  'fs-extra' : fs
  'handlebars'
  'path'
  'require-yaml'
}

class AwsTerraformFileBuilder

  ({@app-config, @exocom-port, @exocom-dns}) ->
    @production-config = @app-config.environments.production
    @terraform-path = '/usr/src/terraform'
    fs.ensure-dir-sync @terraform-path


  generate-terraform: ->
    @_generate-provider-credentials!
    @_generate-main!
    @_generate-vpc!
    @_generate-exocom-cluster!
    @_generate-exocom-service!
    @_generate-cluster 'public'
    @_generate-cluster 'private'
    @_generate-services 'public'
    @_generate-services 'private'


  _generate-main: ->
    # copy cluster creation/destruction policies
    @_copy-template-file 'policies'

    data =
      app-name: @app-config.name
      public-cluster-size: @production-config.providers.aws['public-cluster-size']
      private-cluster-size: @production-config.providers.aws['private-cluster-size']
      domain-name: @production-config.domain

    @_render-template {data, file-name: 'main.tf'}


  _generate-provider-credentials: ->
    data =
      access-key: process.env.AWS_ACCESS_KEY_ID
      secret-key: process.env.AWS_SECRET_ACCESS_KEY
      region: @production-config.providers.aws.region

    @_render-template {data, file-name: 'provider.tf'}


  _generate-vpc: ->
    @_copy-template-file 'vpc'


  _generate-exocom-cluster: ->
    @_copy-template-file 'exocom'


  _generate-exocom-service: ->
    @_build-exocom-container-definition!
    data = {@exocom-dns}

    @_append-to-main-script {data, template-name: 'exocom-service.tf'}


  _generate-cluster: (type) ->
    @_copy-template-file "#{type}-services"


  _generate-services: (type) ->
    for service-name, service-data of @app-config.services["#{type}"]
      service-config = require path.join('/var/app', service-data.location, 'service.yml')
      @_build-service-container-definition service-name, (@_get-image-name service-data), service-config
      data =
        name: service-name
        public-port: service-config.docker['public-port']
        public-url: service-config['deployment-url']
      @_append-to-main-script {data, template-name: "#{type}-service.tf"}


  # renders and writes a Terraform file given a template
  _render-template: ({data, file-name}) ->
    src-path = path.join __dirname, "../../templates/aws-terraform/#{file-name}"

    src = fs.read-file-sync src-path, 'utf-8'
    rendered-file = handlebars.compile(src) data
    fs.write-file-sync "#{@terraform-path}/#{file-name}", rendered-file


  # appends a block of Terraform code to main deployment script
  _append-to-main-script: ({data, template-name}) ->
    src-path = path.join __dirname, "../../templates/aws-terraform/#{template-name}"

    src = fs.read-file-sync src-path, 'utf-8'
    rendered-text = handlebars.compile(src) data
    fs.append-file-sync "#{@terraform-path}/main.tf", rendered-text


  _copy-template-file: (file-name) ->
    src-path = path.join __dirname, "../../templates/aws-terraform/#{file-name}"
    fs.copy-sync src-path, "#{@terraform-path}/#{file-name}"


  _build-service-container-definition: (service-name, image-name, service-config) ->
    container-definition = [
      name: "exosphere-#{service-name}-service"
      image: image-name
      cpu: service-config.docker.cpu
      memory: service-config.docker.memory
      command: service-config.command
      links: service-config.docker.links
      port-mappings: @_build-port-mappings service-config
      environment:
        * name: 'EXOCOM_HOST'
          value: @exocom-dns
        * name: 'EXOCOM_PORT'
          value: "#{@exocom-port}"
        * name: 'SERVICE_NAME'
          value: service-name
      ]
    target-path = path.join @terraform-path, "#{service-name}-container-definition.json"
    fs.write-file-sync target-path, JSON.stringify(container-definition, null, 2)


  # maps public facing port of services to port 80 of the host-port
  # assumes there is only one such public facing port per service
  _build-port-mappings: (service-config) ->
    port-mappings = [
      host-port: 0
      container-port: @exocom-port
      protocol: 'tcp'
    ]
    if service-config.docker['public-port']
      port-mappings.push do
        host-port: 80
        container-port: service-config.docker['public-port']
        protocol: 'tcp'
    port-mappings


  _build-exocom-container-definition: ->
    container-definition = [
      name: 'exocom'
      image: 'originate/exocom'
      cpu: 100
      memory: 500
      essential: true
      command: ['bin/exocom']
      port-mappings: [
        host-port: @exocom-port
        container-port: @exocom-port
        protocol: 'tcp'
      ]
      environment: [
        name: 'SERVICE_MESSAGES'
        value: compile-service-messages(@app-config, '/var/app') |> JSON.stringify
      ]
    ]
    target-path = path.join @terraform-path, 'exocom-container-definition.json'
    fs.write-file-sync target-path, JSON.stringify(container-definition, null, 2)


  _get-image-name: (service-data) ->
    service-config = require path.join('/var/app', service-data.location, 'service.yml')
    "#{service-config.author}/#{service-config.title |> (.replace /\s/g, '-')}"
    #TODO: get image name if location is docker on dockerhub


module.exports = AwsTerraformFileBuilder
