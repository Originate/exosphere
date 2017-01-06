require! {
  'async'
  'path'
  'yaml-cutter'
}


module.exports = ->

  @Given /^I am in the directory of an application with the services:$/, (table) ->
    @app-dir := path.join process.cwd!, 'tmp', 'app'
    @create-empty-app('app')
      .then ~>
        tasks = for row in table.hashes!
          options =
            file: path.join @app-dir, 'application.yml'
            root: 'services.public'
            key: row.NAME
            value: {location: "./#{row.NAME}"}
          yaml-cutter.insert-hash(options, _)
        async.series tasks
      .then ~>
        @write-services table, @app-dir
