require! {
  'chalk' : {black, blue, bold, cyan, dim, green, magenta, red, reset, white, yellow}
  'prelude-ls' : {map, maximum}
}


class Logger

  (roles = []) ->
    @colors =
      exocom: blue
      exorun: reset
      'exo-clone': reset
      'exo-setup': reset
      'exo-test': reset
      'exo-sync': reset
      'exo-lint': reset
      'exo-deploy': reset
    @set-colors roles


  log: ({name, text, trim}) ~>
    color = @colors[name] ? reset
    text = text.trim! if trim
    for line in text.split '\n'
      console.log color(bold "#{@_pad name} "), color(line)

  error: ({name, text, trim}) ~>
    color = @colors[name] ? reset
    text = text.trim! if trim
    for line in text.split '\n'
      console.error color(bold "#{@_pad name} "), red(line)

  # This method may be called after initialization to set/reset colors,
  # given a new list of roles
  set-colors: (roles) ->
    for role, i in roles
      @colors[role] = Logger._default_colors[i % Logger._default_colors.length]
    @length = map (.length), Object.keys(@colors) |> maximum


  @_default_colors = [blue, cyan, magenta, yellow]


  _pad: (text) ->
    "               #{text}".slice -@length



module.exports = Logger
