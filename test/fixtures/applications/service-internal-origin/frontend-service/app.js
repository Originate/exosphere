const express = require('express')
const request = require('request')
const app = express()

app.get('/', (req, res) =>
  request(process.env.BACKEND_ORIGIN, (error, response, body) => {
    res.send(`<html><body>${body}</body></html>`)
  })
)

app.listen(3000, () => console.log("frontend service online"))
