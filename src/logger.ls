require! {
  'chalk' : {black, blue, bold, cyan, dim, green, magenta, red, reset, white, yellow}
  'prelude-ls' : {map, maximum}
}


class Logger

  (service-names) ->
    @colors =
      exocom: blue
      exorun: reset
      'exo-setup': reset
    for service-name, i in service-names
      @colors[service-name] = Logger._default_colors[i]
    @length = map (.length), Object.keys(@colors) |> maximum


  log: ({name, text, trim}) ->
    color = @colors[name]
    text = text.trim! unless trim is false
    for line in text.split '\n'
      console.log color(bold "#{@_pad name} "), color(line)


  @_default_colors = [magenta, cyan, yellow]


  _pad: (text) ->
    "               #{text}".slice -@length



module.exports = Logger
