require! {
  'child_process'
  'dim-console'
  '../../../exosphere-shared' : {templates-path}
  'fs-extra' : fs
  'observable-process' : ObservableProcess
  'path'
  'tmp'
  'tmplconv'
}


#We need to share this variable across scenarios
# for the end-to-end tests
app-dir = null

World = !->

  @create-empty-app = (app-name) ->
    app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join templates-path, 'create-app'
    tmplconv.render(src-path, app-dir, {data})





module.exports = ->
  @app-dir = app-dir
  @World = World
