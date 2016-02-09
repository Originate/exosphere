require! {
  'chalk' : {dim}
}


module.exports =
  log: (text) -> console.log dim text
  error: (text) -> console.log dim text
