require! {
  'nitroglycerin' : N
}


module.exports = (done) ->
  ps-tree process.pid, N (children) ~>
    for child in children
      try
        process.kill child.PID
