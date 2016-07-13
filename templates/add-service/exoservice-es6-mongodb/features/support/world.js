const ejs = require('ejs')


function World() {

  // Fills in _____modelName_____ ids in the placeholders of the template
  this.fillIn_____modelName@camelcase_____Ids = function(template, done) {
    var neededIds = []
    ejs.render(template, {'idOf': (_____modelName_____) => neededIds.push(_____modelName_____) })
    if (neededIds.length === 0) return done(template)
    this.exocom.send({ service: '_____serviceName_____',
                              name: '_____modelName_____.read',
                              payload: {name: neededIds[0]} })
    this.exocom.onReceive( () => {
      const id = this.exocom.receivedMessages[0].payload.id
      done(ejs.render(template, { 'idOf': (_____modelName_____) => id }))
    })
  }


  this.removeIds = function(payload) {
    for (let key in payload) {
      const value = payload[key]
      if (key === 'id' || key === '_id') {
        delete payload[key]
      } else if (typeof value === 'object') {
        payload[key] = this.removeIds(value)
      } else if (typeof value === 'array') {
        const temp = []
        for (const child in value) {
          temp.push(this.removeIds(child))
        }
        payload[key] = temp
      }
    }
    return payload
  }

}


module.exports = function() {
  this.World = World
}
