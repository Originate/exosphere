require! {
  \prelude-ls : {any, map}
  'dockerode' : Docker
  'wait' : {wait}
}

docker = new Docker

# Helper class used to manage Docker processes not started by docker-compose
class DockerHelper

  @start-container = ({Image, name, HostConfig}, done) ->
    DockerHelper.list-running-containers (err, running-containers) ->
      | err                              => done err
      | running-containers.includes name => done!
      | otherwise                        =>
        docker.create-container {Image, name, HostConfig}, (err, container) -> 
          | err => done err
          container.start (err, data) ->
            | err => done err
            wait 2_000, done


  @remove-container = ({name}, done) ->
    docker.list-containers {name}, (err, containers) ->
      | err => done err
      DockerHelper._force-remove-containers containers, done


  @remove-all-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err                => done err
      | !containers.length => done!
      | otherwise          => DockerHelper._force-remove-containers containers, done


  @list-running-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err => done err
      # Names field is printed like: Names: [ '/exocom' ]
      done null, map((.Names?[0] |> (.replace '/', '')), containers)


  @pull-image = (image, done) ->
    docker.pull image, (err, stream) ->
      | err => done err
      done!


  @has-image = (image, done) ->
    docker.list-images (err, images) ->
      | err => done err
      done null, any((.RepoTags?[0].includes image), images)


  @list-images = (done) ->
    docker.list-images (err, images) ->
      | err => done err
      done null, map((.RepoTags?[0]), images)


  @_force-remove-containers = (containers, done) ->
    for container in containers
      docker.get-container(container.Id).remove {force:true}, (err) ->
        | err => done err
        done!


module.exports = DockerHelper
