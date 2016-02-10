require! {
  'chalk' : {blue, bold, cyan, dim, green, magenta, red, white}
}


class Logger

  ->
    @colors =
      exocomm: blue
      users: magenta
      web: cyan


  log: ({name, text}) ->
    color = @colors[name]
    console.log color(bold " #{@_pad name} "), color(text.trim!)


  _pad: (text) ->
    "     #{text}".slice -7


module.exports = Logger
