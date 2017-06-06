require! {
  'chalk' : {cyan, green, red, blue}
  'child_process': child-process
}

module.exports = ->

  if process.argv[2] is "help"
    return help!

  console.log 'We are about to clean up your Docker workspace!\n'

  get-dangling-images-cmd = 'docker images -f dangling=true -q'
  get-dangling-volumes-cmd = 'docker volume ls -qf dangling=true'
  remove-images-cmd = 'docker rmi'
  remove-volumes-cmd = 'docker volume rm'
  remove-images-msg = "\nremoved all dangling images"
  remove-volumes-msg = "\nremoved all dangling volumes"

  remove remove-volumes-cmd, get-dangling-volumes-cmd, remove-images-msg
  remove remove-images-cmd, get-dangling-images-cmd, remove-volumes-msg


function help
  help-message =
    """
    \nUsage: #{cyan 'exo clean'}

    Removes dangling Docker images and volumes.
    """
  console.log help-message


function remove (remove-cmd, get-ids-cmd, success-msg)
  child-process.exec(get-ids-cmd, (error, stdout, stderr) ->
    if error
      console.error red error;
      return
    else if stdout
      cmd = "#remove-cmd #stdout".replace(/[\r\n]/g ' ')
      child-process.exec(cmd, (error, stdout, stderr) ->
        if error
          console.error error;
          return
      )
    console.log green success-msg
  )
