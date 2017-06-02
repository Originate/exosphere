require! {
  \observable-process : ObservableProcess
}


class DockerCompose

  @build-all-images = ({write, cwd}, done) ->
    new ObservableProcess(['docker-compose', 'build']
                          cwd: cwd
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done

  @pull-all-images = ({write, cwd}, done) ->
    new ObservableProcess(['docker-compose', 'pull']
                          cwd: cwd
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @run-all-images = ({env, cwd, write}, done) ->
    new ObservableProcess(['docker-compose', 'up']
                          cwd: cwd
                          env: env
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @kill-container = ({service-name, cwd, write}, done) ->
    new ObservableProcess(['docker-compose', 'kill', 'service-name']
                          cwd: cwd
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @kill-all-containers = ({write, cwd}, done) ->
    new ObservableProcess(['docker-compose', 'down']
                          cwd: cwd
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @create-new-container = ({service-name, cwd, env, write}, done) ->
    new ObservableProcess(['docker-compose', 'create', '--build', 'service-name']
                          cwd: cwd
                          env: env
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


  @start-container = ({service-name, cwd, env, write}, done) ->
    new ObservableProcess(['docker-compose', 'restart', 'service-name']
                          cwd: cwd
                          env: env
                          stdout: {write}
                          stderr: {write})
      ..on 'ended', done


module.exports = DockerCompose
