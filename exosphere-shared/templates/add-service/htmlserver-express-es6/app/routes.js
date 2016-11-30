module.exports = ({GET, resources}) => {

  GET('/', { to: 'index#index' })

  // resources('users', { only: ['create', 'destroy'] })

}
