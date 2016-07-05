// Middleware to compile assets on the fly in development mode
// In production we use the asset_compiler to pre-compile everything to public/assets
const webpack = require('webpack')
const webpackDevMiddleware = require('webpack-dev-middleware')
const webpackConfig = require('./webpack-config')

module.exports = webpackDevMiddleware(webpack(webpackConfig),
                                      { publicPath: '/assets/',
                                        noInfo: true })
