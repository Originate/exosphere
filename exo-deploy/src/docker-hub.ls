require! {
  'async'
  'child_process'
  'dashify'
  '../../exosphere-shared' : {DockerHelper}
  'fs'
  'js-yaml' : yaml
  'observable-process' : ObservableProcess
  'path'
}


# Pushes images to DockerHub
class DockerHub

  (@app-config, @logger) ->


  push: (done) ->
    images = @_image-names!
    for image in images
      if !DockerHelper.image-exists image then return done new Error "No Docker image exists for service '#{image.name}'. Please run exo-setup."

    async.each-series images, (~> @_push-image &0, &1), done


  _push-image: (image, done) ~>
    @logger.log role: 'exo-deploy', text: "pushing #{image.name} to DockerHub..."
    new ObservableProcess("docker push #{image.author}/#{image.name}",
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code  =>  return done new Error("#{image.name} could not be pushed to DockerHub")
        @logger.log role: 'exo-deploy', text: "#{image.name} pushed to DockerHub"
        done!


  _image-names: ->
    names = []
    for service-type of @app-config.services
      for name, config of @app-config.services[service-type]
        service-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, config.location, 'service.yml'), 'utf8')
        names.push do
          author: service-config.author
          name: dashify service-config.type
          #TODO: get image name if location is docker on dockerhub
    names


  write: (text) ~>
    @logger.log {role: 'exo-deploy', text, trim: no}


module.exports = DockerHub
