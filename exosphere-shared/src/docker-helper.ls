require! {
  \prelude-ls : {any, head, map}
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


  @remove-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err                => done err
      | !containers.length => done!
      | otherwise          => DockerHelper._force-remove-containers containers, done


  @list-running-containers = (done) ->
    docker.list-containers (err, containers) ->
      | err => done err
      # Names field is printed like: Names: [ '/exocom' ]
      done null, map((.Names?[0] |> (.replace '/', '')), containers)


  @list-images = (done) ->
    docker.list-images (err, images) ->
      | err => done err
      # Image name is printed like: RepoTags: [ 'exocom:latest' ]
      done null, map((.RepoTags |> head |> (.split ':') |> head), images)


  @_force-remove-containers = (containers, done) ->
    for container in containers
      docker.get-container(container.Id).remove {force:true}, (err) ->
        | err => done err
        done!


module.exports = DockerHelper
