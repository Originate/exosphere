module.exports = ({GET, resources}) ->

  GET '/' to: 'home#index'

  resources 'tweets' only: <[ create destroy ]>
  resources 'users'
