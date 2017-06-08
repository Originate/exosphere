require! {
  'chalk' : {black, blue, bold, cyan, dim, green, magenta, red, reset, white, yellow, strip-color, has-color}
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


  log: ({role, text, trim}) ~>
    color = @colors[role] ? reset
    text = text.trim! if trim
    for line in text.split '\n'
      @_parse-line role, line, (left, right) ~>
        console.log color(bold "#{@_pad left} "), color(right)


  error: ({role, text, trim}) ~>
    color = @colors[role] ? reset
    text = text.trim! if trim
    for line in text.split '\n'
      @_parse-line role, line, (left, right) ~>
        console.error color(bold "#{@_pad left} "), red(right)


  # This method may be called after initialization to set/reset colors,
  # given a new list of roles
  set-colors: (roles) ->
    for role, i in roles
      @colors[role] = Logger._default_colors[i % Logger._default_colors.length]
    @length = map (.length), Object.keys(@colors) |> maximum


  @_default_colors = [blue, cyan, magenta, yellow]


  _parse-line: (role, line, done) ->
    segments = [segment.trim! for segment in line / /\s+\|\s*/]
    if segments.length == 2
      service = @_parse-service segments[0]
      done service, (@_reformat-line(segments[1]))
    else
      done role, line


  _parse-service: (text) ->
    text - /(\d+\.)?(\d+\.)?(\*|\d+)$/


  _reformat-line: (line) ->
    color-str = @_get-color-str line
    "#color-str#{(strip-color line).trim!}"  


  _get-color-str: (styled-string) ->
    color-strings = (/\x1b[^m]*m/).exec styled-string
    if color-strings then color-strings[0] else ''


  _pad: (text) ->
    color-str = @_get-color-str text
    "#color-str#{"               #{strip-color text}".slice -@length}"


module.exports = Logger
