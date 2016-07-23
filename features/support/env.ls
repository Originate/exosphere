require! {
  'fs-extra' : fs
  'nitroglycerin' : N
  'ps-tree'
}


module.exports = ->

  @set-default-timeout 2000


  @After tags: ['~@e2e'], (scenario, done) ->
    if @app-dir
      fs.remove-sync @app-dir
    ps-tree process.pid, N (children) ~>
      for child in children
        try
          process.kill child.PID
      done!
