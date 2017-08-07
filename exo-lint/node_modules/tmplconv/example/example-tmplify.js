'use strict'

const tmplconv = require('tmplconv')

// Generate template from existing directory
tmplconv.tmplify('demo/demo-app', 'asset/app-tmpl', {
  // Patterns of files to tmplify
  pattern: [
    'lib/*.js',
    'test/*_test.js'
  ],
  // Rule to tmplify
  data: {
    'name': 'my-awesome-app',
    'description': "This is an example for the app templates."
  }
}).then((result) => {
  /* ... */
})
