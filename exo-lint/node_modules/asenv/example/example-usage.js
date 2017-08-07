'use strict'

const { getEnv, setEnv, isProduction } = require('asenv')

{
  let env = getEnv()
  console.log('env=', env)

  /* ... */

  setEnv('production')

  /* ... */

  if (isProduction()) {
    /* ... */
  }
}
