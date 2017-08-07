'use strict'

const ababel = require('ababel')
const co = require('co')

co(function * () {
  yield ababel('**/*.jsx', {
    cwd: 'src',
    out: 'dest',
    minified: true,
    presets: [ 'es2015', 'react' ]
  })
}).catch((err) => console.error(err))

