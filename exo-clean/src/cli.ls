require! {
  'chalk' : {cyan, green, red, blue}
  '../../exosphere-shared' : {DockerHelper}
}

module.exports = ->

  if process.argv[2] is "help"
    return help!

  console.log 'We are about to clean up your Docker workspace!\n'

  DockerHelper.get-dangling-images (err, images) ->
    | err => throw err
    DockerHelper.force-remove-images images, (err) ->
      | err => throw err
      console.log green 'removed all dangling images'
  DockerHelper.get-dangling-volumes (err, volumes) ->
    | err => throw err
    DockerHelper.force-remove-volumes volumes, (err) ->
      | err => throw err
      console.log green 'removed all dangling volumes'


function help
  help-message =
    """
    \nUsage: #{cyan 'exo clean'}

    Removes dangling Docker images and volumes.
    """
  console.log help-message
