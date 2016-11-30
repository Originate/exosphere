require! {
  'nitroglycerin' : N
  'ps-tree'
}


module.exports = (done) ->
  ps-tree process.pid, N (children) ~>
    for child in children
      try
        process.kill child.PID
    done!
