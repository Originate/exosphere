require! {
  'fs-extra' : fs
  'path'
  'tmplconv'
}



World = !->

  @create-empty-app = (app-name) ->
    @app-dir = path.join process.cwd!, 'tmp', app-name
    fs.empty-dir-sync @app-dir
    data =
      'app-name': app-name
      'app-description': 'Empty test application'
      'app-version': '1.0.0'
    src-path = path.join process.cwd!, '..', 'exosphere-shared', 'templates', 'create-app'
    tmplconv.render(src-path, @app-dir, {data, silent: true})


  @write-services = (table, @app-dir) ->
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
      fs.mkdir-sync path.join(@app-dir, row.NAME)
      fs.write-file-sync path.join(@app-dir, row.NAME, 'service.yml'), content



module.exports = World
