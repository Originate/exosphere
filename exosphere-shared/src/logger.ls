require! {
  'chalk' : {black, blue, bold, cyan, dim, green, magenta, red, reset, white, yellow, strip-color}
  'pad-left'
  'prelude-ls' : {map, maximum}
}


class Logger

  (roles = []) ->
    @colors =
      exocom: cyan
      exorun: reset
      'exo-clone': reset
      'exo-setup': reset
      'exo-test': reset
      'exo-sync': reset
      'exo-lint': reset
      'exo-deploy': reset
    @set-colors roles


  log: ({role, text, trim}) ~>
    text = text.trim! if trim
    for line in text.split '\n'
      parsed-line = @_parse-line role, line
      color = @colors[parsed-line.left] ? reset
      console.log color("#{bold "#{@_pad "#{parsed-line.left}"} "} #{parsed-line.right}")


  error: ({role, text, trim}) ~>
    text = text.trim! if trim
    for line in text.split '\n'
      parsed-line = @_parse-line role, line
      color = @colors[parsed-line.left] ? reset
      console.error color(bold "#{@_pad "#{parsed-line.left}"} "), red(parsed-line.right)


  # This method may be called after initialization to set/reset colors,
  # given a new list of roles
  set-colors: (roles) ->
    for role, i in roles
      @colors[role] = Logger._default_colors[i % Logger._default_colors.length]
    @length = map (.length), Object.keys(@colors) |> maximum


  @_default_colors = [magenta, blue, yellow, cyan]


  _parse-line: (role, line) ->
    segments = [segment.trim! for segment in line / /\s+\|\s*/]
    if segments.length == 2
      service = @_parse-service segments[0]
      {left: service, right: (@_reformat-line(segments[1]))}
    else
      {left:role, right: line}


  _parse-service: (text) ->
    strip-color text - /(\d+\.)?(\d+\.)?(\*|\d+)$/


  _reformat-line: (line) ->
    "#{(strip-color line).trim!}"


  _pad: (text) ->
    pad-left(text, @length, ' ')


module.exports = Logger
