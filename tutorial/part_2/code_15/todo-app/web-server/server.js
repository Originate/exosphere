const ExoRelay = require('exorelay');
const express = require('express');

const exoRelay = new ExoRelay({serviceName: process.env.SERVICE_NAME,
                               exocomPort: process.env.EXOCOM_PORT});
exoRelay.listen(process.env.EXORELAY_PORT);
exoRelay.on('error', (err) => { console.log(`Error: %{err}`); });
exoRelay.on('online', () => {

  const app = express();
  app.set('view engine', 'jade');

  app.get('/', (req, res) => {
    exoRelay.send('todos.list', (todos) => {
      res.render('index', {todos: todos});
    });
  });

  app.listen(3000, () => {
    console.log('Todo web server listening on port 3000');
  });

});
