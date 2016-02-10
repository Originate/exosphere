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
    console.log color(bold " #{name} "), color(text.trim!)



module.exports = Logger
