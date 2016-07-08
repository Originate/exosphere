module.exports = ({GET, resources}) => {

  GET('/', { to: 'index#index' })

  resources('tweets', { only: ['create', 'destroy'] })
  resources('users')

}
