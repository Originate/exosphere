'use strict'

const writeout = require('writeout')

// Generate a file.
writeout('hello-wold', 'This is the contents text', {
  mkdirp: true,
  skipIfIdentical: true
}).then((result) => {
  if (!result.skipped) {
    console.log('File generated:', result.filename)
  }
}).catch((err) =>
  console.error(err)
)

