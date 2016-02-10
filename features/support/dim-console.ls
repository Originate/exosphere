require! {
  'chalk' : {dim}
}


module.exports =
  log: (text) -> console.log dim text.replace /\n*$/, ''
  error: (text) -> console.log dim text.replace /\n*$/, ''
