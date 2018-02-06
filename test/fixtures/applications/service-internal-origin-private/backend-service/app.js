const express = require('express')
const app = express()

app.get('/', (req, res) => res.send('Backend service content'))

app.listen(4000, () => console.log("backend service online"))
