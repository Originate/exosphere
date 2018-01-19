const express = require('express')
const app = express()

app.get('/', (req, res) => res.send('Backend service content'))

app.listen(8080, () => console.log("backend service online"))
