#!/usr/bin/env/node
'use strict'

const aglob = require('aglob')

aglob([
  'lib/*.js',
  'doc/**/.js'
], {
  cwd: process.cwd(),
  ignore: []
}).then((filenames) => {
  console.log(filenames)
})
