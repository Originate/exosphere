require! {
  \observable-process : ObservableProcess
}


class DockerCompose

  @build-all-images = ({write}, done) ->
    new ObservableProcess('docker-compose build'
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done

  @pull-all-images = ({write}, done) ->
    new ObservableProcess('docker-compose pull'
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @run-all-images = ({env, write}, done) ->
    new ObservableProcess('docker-compose up'
                          env: env
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @kill-container = ({service-name, write}, done) ->
    new ObservableProcess("docker-compose kill #{service-name}"
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @kill-all-containers = ({write}, done) ->
    new ObservableProcess('docker-compose down'
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @create-new-container = ({service-name, env, write}, done) ->
    new ObservableProcess("docker-compose create --build #{service-name}"
                          env: env
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @start-container = ({service-name, env, write}, done) ->
    new ObservableProcess("docker-compose restart #{service-name}"
                          env: env
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


module.exports = DockerCompose
