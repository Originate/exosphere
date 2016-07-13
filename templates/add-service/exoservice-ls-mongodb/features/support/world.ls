require! {
  'eco'
}


World = !->

  # Fills in _____modelName_____ ids in the placeholders of the template
  @fill-in-_____modelName_____-ids = (template, done) ->
    needed-ids = []
    eco.render template, id_of: (_____modelName_____) -> needed-ids.push _____modelName_____
    return done template if needed-ids.length is 0
    @exocom
      ..send service: '_____serviceName_____', name: '_____modelName_____.read', payload: {name: needed-ids[0]}
      ..on-receive ~>
        id = @exocom.received-messages[0].payload.id
        done eco.render(template, id_of: (_____modelName_____) -> id)


  @remove-ids = (payload) ->
    for key, value of payload
      if key is 'id'
        delete payload[key]
      else if typeof value is 'object'
        payload[key] = @remove-ids value
      else if typeof value is 'array'
        payload[key] = [@remove-ids(child) for child in value]
    payload



module.exports = ->
  @World = World
