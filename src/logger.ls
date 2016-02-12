require! {
  'chalk' : {black, blue, bold, cyan, dim, green, magenta, red, white, yellow}
  'prelude-ls' : {map, maximum}
}


class Logger

  (service-names) ->
    @length = map (.length), service-names |> maximum
    @colors =
      exocomm: blue
      exorun: -> it   # use the default color here
      users: magenta
      web: cyan
      dashboard: yellow


  log: ({name, text}) ->
    color = @colors[name]
    console.log color(bold "#{@_pad name} "), color(text.trim!)


  _pad: (text) ->
    "     #{text}".slice -@length


module.exports = Logger
