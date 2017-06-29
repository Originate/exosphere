require! {
  'prelude-ls' : {map, sum}
}


module.exports =

  ascii: (text) ->
    [0 til text.length]
    |> map text~char-code-at
    |> sum
