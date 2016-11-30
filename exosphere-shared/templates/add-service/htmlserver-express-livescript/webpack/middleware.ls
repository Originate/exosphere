# Middleware to compile assets on the fly in development mode
# In production we use the asset_compiler to pre-compile everything to public/assets
require! {
  'webpack'
  'webpack-dev-middleware'
  './webpack-config'
}


module.exports = webpack-config
  |> webpack
  |> webpack-dev-middleware _,
    public-path: '/assets/'
    no-info: yes
