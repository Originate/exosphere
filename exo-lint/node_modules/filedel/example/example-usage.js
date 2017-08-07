'use strict'

const filedel = require('filedel')

// Generate a file.
filedel('/src/*.tmp', {
  force: true
}).then(() => {
  /* ... */
})
