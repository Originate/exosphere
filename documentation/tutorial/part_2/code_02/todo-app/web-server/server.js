const express = require('express');
const app = express();
app.set('view engine', 'jade');

app.get('/', (req, res) => {
  res.render('index');
});

app.listen(3000, () => {
  console.log('Todo web server listening on port 3000');
});
