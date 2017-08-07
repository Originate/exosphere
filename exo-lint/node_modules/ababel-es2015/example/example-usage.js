'use strict'

const ababelEs2015 = require('ababel-es2015')

const co = require('co')

co(function * () {
  yield ababelEs2015('**/*.js', {
    cwd: 'src',
    out: 'dest',
    minified: true
  })
}).catch((err) => console.error(err))