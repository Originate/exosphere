require! {
  'async'
  'child_process'
  '../../../exosphere-shared' : {DockerHelper}
  'observable-process' : ObservableProcess
  'path'
  'require-yaml'
}


# Pushes images to DockerHub
class DockerHub

  (@app-config, @logger) ->


  push: (done) -> # TODO: make sure this is run before bin/start-deploy
    images = @_image-names!
    for image in images
      if !DockerHelper.image-exists image then return done new Error "No Docker image exists for service '#{image.name}'. Please run exo-setup."

    async.each-series images, (~> @_push-image &0, &1), done


  _push-image: (image, done) ~>
    @logger.log name: 'exo-deploy', text: "pushing #{image.name} to DockerHub..."
    new ObservableProcess("docker push #{image.author}/#{image.name}",
                          stdout: {@write}
                          stderr: {@write})
      ..on 'ended', (exit-code) ~>
        | exit-code  =>  return done new Error("#{image.name} could not be pushed to DockerHub")
        @logger.log name: 'exo-deploy', text: "#{image.name} pushed to DockerHub"
        done!


  _image-names: ->
    names = []
    for service-type of @app-config.services
      for name, config of @app-config.services[service-type]
        service-config = require path.join(process.cwd!, config.location, 'service.yml')
        names.push do
          author: service-config.author
          name: path.basename config.location
    names


  write: (text) ~>
    @logger.log {name: 'exo-deploy', text, trim: no}


module.exports = DockerHub
