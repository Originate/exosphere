require! {
  'chalk' : {black, blue, bold, cyan, dim, green, magenta, red, reset, white, yellow, strip-color}
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
      @_log-line role, line, color

  error: ({role, text, trim}) ~>
    color = @colors[role] ? reset
    text = text.trim! if trim
    for line in text.split '\n'
      @_log-line line

  # This method may be called after initialization to set/reset colors,
  # given a new list of roles
  set-colors: (roles) ->
    for role, i in roles
      @colors[role] = Logger._default_colors[i % Logger._default_colors.length]
    @length = map (.length), Object.keys(@colors) |> maximum


  @_default_colors = [blue, cyan, magenta, yellow]


  _pad: (text, offset=0) ->
    "               #{text}".slice -@length


  _log-line: (role, line, color) ->
    elts = [elt.trim! for elt in line / '|']
    if elts.length == 2
      color-str = @_get-color-str elts[0]
      service = @_parse-service elts[0]
      console.log bold color-str + "#{@_pad service, color-str.length} ", elts[1]
    else
      console.log color(bold "#{@_pad role} "), color(line)


  _get-color-str: (styled-string) ->
    color-strings = (/\x1b[^m]*m/).exec styled-string
    if color-strings then color-strings[0] else ''


  _parse-service: (text) ->
    strip-color text - /(\d+\.)?(\d+\.)?(\*|\d+)$/


module.exports = Logger
