require! {
  \prelude-ls : {any}
  'dockerode' : Docker
  'wait' : {wait}
}

docker = new Docker

# Helper class used to manage Docker processes not started by docker-compose
class DockerHelper

  @start-container = ({Image, name, HostConfig}, done) ->
    DockerHelper._container-is-running name, (err, is-running) ->
      | err => done err
      | is-running => done!
      | otherwise =>
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
      | err => done err
      DockerHelper._force-remove-containers containers, done


  @_force-remove-containers = (containers, done) ->
    for container in containers
      docker.get-container(container.Id).remove {force:true}, (err) ->
        | err => done err
        done!


  @_container-is-running = (name, done) ->
    docker.list-containers {name}, (err, containers) ->
      | err => done err
      done null, any (.Names.includes name), containers


module.exports = DockerHelper
