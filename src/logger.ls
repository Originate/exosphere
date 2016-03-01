require! {
  'chalk' : {black, blue, bold, cyan, dim, green, magenta, red, white, yellow}
  'prelude-ls' : {map, maximum}
}


class Logger

  (service-names) ->
    @colors =
      exocomm: blue
      exorun: -> it   # use the default color here
      'exo-install': -> it   # use the default color here
    for service-name, i in service-names
      @colors[service-name] = Logger._colors[i]
    @length = map (.length), Object.keys(@colors) |> maximum


  log: ({name, text}) ->
    color = @colors[name]
    console.log color(bold "#{@_pad name} "), color(text.trim!)


  @_colors = [magenta, cyan, yellow]


  _pad: (text) ->
    "               #{text}".slice -@length



module.exports = Logger
