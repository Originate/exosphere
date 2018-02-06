const express = require('express')
const request = require('request')
const app = express()

app.get('/', (req, res) =>
  request(process.env.BACKEND_INTERNAL_ORIGIN, (error, response, body) => {
    res.send(`<html><body>${body}</body></html>`)
  })
)

app.listen(5000, () => console.log("frontend service online"))
