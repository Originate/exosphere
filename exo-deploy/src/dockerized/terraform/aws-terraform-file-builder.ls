require! {
  'async'
  'dashify'
  'fs-extra' : fs
  'handlebars'
  'js-yaml' : yaml
  'path'
}

class AwsTerraformFileBuilder

  ({@app-config, @exocom-port, @exocom-dns}) ->
    @production-config = @app-config.environments.production
    @terraform-path = '/usr/src/terraform'
    fs.ensure-dir-sync @terraform-path


  empty-terraform-dir: ->
    fs.empty-dir-sync '/usr/src/terraform', (err) ->
      | err => process.stdout.write "Could not clear directory '/usr/src/terraform' inside Docker #{err.message}"


  generate-terraform: ->
    @generate-provider-credentials!
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


  generate-provider-credentials: ->
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
      service-config = yaml.safe-load fs.read-file-sync(path.join('/var/app', service-data.location, 'service.yml'), 'utf8')
      @_build-service-container-definition service-name, (@_get-image-name service-data), service-config
      data =
        name: service-name
        public-port: service-config.production.aws['public-port']
        public-url: service-config.production.url
      @_append-to-main-script {data, template-name: "#{type}-service.tf"}


  # renders and writes a Terraform file given a template
  _render-template: ({data, file-name}) ->
    src-path = path.join __dirname, "../../../templates/aws-terraform/#{file-name}"

    src = fs.read-file-sync src-path, 'utf-8'
    rendered-file = handlebars.compile(src) data
    fs.write-file-sync "#{@terraform-path}/#{file-name}", rendered-file


  # appends a block of Terraform code to main deployment script
  _append-to-main-script: ({data, template-name}) ->
    src-path = path.join __dirname, "../../../templates/aws-terraform/#{template-name}"

    src = fs.read-file-sync src-path, 'utf-8'
    rendered-text = handlebars.compile(src) data
    fs.append-file-sync "#{@terraform-path}/main.tf", rendered-text


  _copy-template-file: (file-name) ->
    src-path = path.join __dirname, "../../../templates/aws-terraform/#{file-name}"
    fs.copy-sync src-path, "#{@terraform-path}/#{file-name}"


  _build-service-container-definition: (service-name, image-name, service-config) ->
    container-definition = [
      name: "exosphere-#{service-name}-service"
      image: image-name
      cpu: service-config.production.aws.cpu
      memory: service-config.production.aws.memory
      command: service-config.startup.command |> (.split ' ')
      port-mappings: @_build-port-mappings service-config
      environment: @_build-service-environment-variables service-name, service-config
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
    if service-config.production.aws['public-port']
      port-mappings.push do
        host-port: 80
        container-port: service-config.production.aws['public-port']
        protocol: 'tcp'
    port-mappings


  _build-service-environment-variables: (service-name, service-config) ->
    environment =
      * name: 'EXOCOM_HOST'
        value: @exocom-dns
      * name: 'EXOCOM_PORT'
        value: "#{@exocom-port}"
      * name: 'ROLE'
        value: service-name
    for dependency of service-config.dependencies
      switch dependency
        | 'mongo' =>
          process.env.MONGODB_USER ? throw new Error "MONGODB_USER not provided"
          process.env.MONGODB_PW ? throw new Error "MONGODB_PW not provided"
          environment ++= {name: 'MONGODB_USER', value: process.env.MONGODB_USER}
          environment ++= {name: 'MONGODB_PW', value: process.env.MONGODB_PW}
    environment


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
        name: 'SERVICE_ROUTES'
        value: @_compile-service-messages! |> JSON.stringify
      ]
    ]
    target-path = path.join @terraform-path, 'exocom-container-definition.json'
    fs.write-file-sync target-path, JSON.stringify(container-definition, null, 2)


  _get-image-name: (service-data) ->
    service-config = yaml.safe-load fs.read-file-sync(path.join('/var/app', service-data.location, 'service.yml'), 'utf8')
    "#{service-config.author}/#{dashify service-config.title}"
    #TODO: get image name if location is docker on dockerhub


  _compile-service-messages: ->
    service-messages = []
    for type of @app-config.services
      for service-name, service-data of @app-config.services["#{type}"]
        service-config = yaml.safe-load fs.read-file-sync(path.join('/var/app', service-data.location, 'service.yml'), 'utf8')
        service-messages.push do
          name: service-name
          receives: service-config.messages.receives
          sends: service-config.messages.sends
          namespace: service-data.namespace
    service-messages


module.exports = AwsTerraformFileBuilder
