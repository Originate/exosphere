require! {
  'child_process'
  'dim-console'
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'tmp'
  'tmplconv'
}


World = !->

  @checkout-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    fs.copy-sync path.join(process.cwd!, 'example-apps', app-name), app-dir


  @setup-app = (app-name, done) ->
    @process = new ObservableProcess(path.join(process.cwd!, 'bin', 'exo setup'),
                                     cwd: path.join(process.cwd!, 'tmp', app-name),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..on 'ended', done


  @start-app = (app-name, done) ->
    @process = new ObservableProcess(path.join('..', 'bin', 'exo run'),
                                     cwd: path.join(process.cwd!, 'tmp'),
                                     stdout: dim-console.process.stdout
                                     stderr: dim-console.process.stderr)
      ..wait "application ready", done

  @create-empty-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join process.cwd!, 'templates', 'create-app'
    tmplconv.render(src-path, app-dir, {data})


  @write-services = (table, app-dir) ->
    for row in table.hashes!
      content = """
        name: #{row.NAME}
        decription: test service

        messages:
        """
      if row.SENDS
        content += "\n sends: "
        for message in row.SENDS.split(', ')
          content += "\n    - #{message}"
      if row.RECEIVES
        content += "\n receives: "
        for message in row.RECEIVES.split(', ')
          content += "\n    - #{message}"
      fs.mkdir-sync path.join(app-dir, row.NAME)
      fs.write-file-sync path.join(app-dir, row.NAME, 'service.yml'), content


  @make-repo = (cwd) ->
    child_process.exec-sync("git init && git add --all && git commit -m \"initial commit\"",
                            cwd: cwd,
                            stdio: [1,2])



module.exports = ->
  @World = World
