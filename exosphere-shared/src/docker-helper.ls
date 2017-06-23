require! {
  'async'
  \prelude-ls : {any, head, map}
  'dockerode' : Docker
  'stream'
  'text-stream-accumulator' : TextStreamAccumulator
  'text-stream-search' : TextStreamSearch
}

docker = new Docker

# Helper class used to manage Docker processes not started by docker-compose
class DockerHelper

  @start-container = ({Image, name, HostConfig, online-text}, done) ->
    DockerHelper.list-running-containers (err, running-containers) ->
      | err                              => done err
      | running-containers.includes name => done!
      | otherwise                        =>
        docker.create-container {Image, name, HostConfig}, (err, container) ->
          | err => done err
          container.attach {stream: true, stdout: true, stederr: true}, (err, stream) ->
            text-stream-search = new TextStreamSearch stream
            container.start (err) ->
              | err => done err
              text-stream-search.wait online-text, done


  @remove-container = ({name}, done) ->
    docker.list-containers {name}, (err, containers) ->
      | err       => done err
      | otherwise => DockerHelper._force-remove-containers containers, done


  @remove-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err                => done err
      | !containers.length => done!
      | otherwise          => DockerHelper._force-remove-containers containers, done


  @list-running-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err       => done err
      | otherwise =>
        # Names field is printed like: Names: [ '/exocom' ]
        done null, map((.Names?[0] |> (.replace '/', '')), containers)


  @list-images = (done) ->
    docker.list-images (err, images) ->
      | err       => done err
      | otherwise =>
        # Image name is printed like: RepoTags: [ 'exocom:latest' ]
        done null, map((.RepoTags |> head |> (.split ':') |> head), images)


  @pull-image = ({image}, done) ->
    docker.pull image, (err, stream) ->
      console.log "Downloading docker image for '#{image}'..."
      docker.modem.follow-progress stream, (err, output) ->
        | err       => done err
        | otherwise => console.log "'#{image}' download complete"; done!


  # Runs cat on file in Docker container to print its content
  @cat-file = ({image, file-name}, done) ->
    DockerHelper.pull-image {image}, (err) ->
      | err       => done err
      | otherwise =>
        stdout-stream = new stream.PassThrough
        text-stream-search = new TextStreamAccumulator stdout-stream
        docker.run image, ['cat', file-name], stdout-stream, (err, data, container) ->
          | err       => done err
          | otherwise => done null, text-stream-search.to-string!


  @pull-image = ({image}, done) ->
    docker.pull image, (err, stream) ->
      console.log "Downloading docker image for '#{image}'..."
      docker.modem.follow-progress stream, (err, output) ->
        | err       => done err
        | otherwise => console.log "'#{image}' download complete"; done!


  @get-dangling-images = (done) ->
    docker.list-images {"filters": '{"dangling": ["true"]}'}, done


  @force-remove-images = (images, done) ->
    async.map-series images, (-> docker.get-image(&0.Id).remove {force:true}, &1), done


  @get-dangling-volumes = (done) ->
    docker.list-volumes {"filters": '{"dangling": ["true"]}'}, (err, volumes) ->
      | err       => done err
      | otherwise => done null, (volumes.Volumes or [])


  @force-remove-volumes = (volumes, done) ->
    async.map-series volumes, (-> docker.get-volume(&0.Name).remove {force:true}, &1), done


  @_force-remove-containers = (containers, done) ->
    async.map-series containers, (-> docker.get-container(&0.Id).remove {force:true}, &1), done


module.exports = DockerHelper
