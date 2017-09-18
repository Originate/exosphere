const path = require('path');

module.exports = {
  entry: ["./src/index.js"],

  output: {
    publicPath: "/",
    path: path.join(__dirname, "dist"),
    filename: "bundle.js"
  },
}
