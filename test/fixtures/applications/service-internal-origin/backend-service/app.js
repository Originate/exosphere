const express = require('express')
const app = express()

app.get('/', (req, res) => res.send('Backend service reached'))

app.listen(3000)
