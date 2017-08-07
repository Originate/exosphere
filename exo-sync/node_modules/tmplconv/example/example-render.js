'use strict'

const tmplconv = require('tmplconv')

// Render files from existing template
tmplconv.render('asset/app-tmpl', 'demo/demo-app', {
  // Data to render
  data: {
    'name': 'my-awesome-app',
    'description': "This is an example for the app templates."
  }
}).then((result) => {
  /* ... */
})
