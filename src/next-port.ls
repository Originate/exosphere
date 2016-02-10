require! {
  'nitroglycerin' : N
  'portfinder'
}


last-port = 5000


next-port = (done) ->
  portfinder
    ..base-port = last-port
    ..get-port N (port) ->
      last-port := port + 1
      done port



module.exports = next-port
