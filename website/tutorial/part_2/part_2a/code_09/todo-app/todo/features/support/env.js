process.env.NODE_ENV = 'test'
const defineSupportCode = require('cucumber').defineSupportCode,
      {MongoClient} = require('mongodb'),
      N = require('nitroglycerin'),
      World = require('./world')


var db = null
const getDb = (done) => {
  if (db) return done(db)
  MongoClient.connect(`mongodb://${process.env.MONGO}:27017/exosphere-todo-test`, N( (mongoDb) => {
    db = mongoDb
    done(db)
  }))
}


defineSupportCode(function({Before, After, AfterAll, setDefaultTimeout, setWorldConstructor}) {

  setDefaultTimeout(1000)
  setWorldConstructor(World)


  Before(function(_scenario, done) {
    getDb( (db) => {
      db.collection('todos').drop(function(err) {
        // ignore errors here, since we are only cleaning up the test database
        // and it might not even exist
        done()
      })
    })
  })


  After(function(_scenario, done) {
    this.process.kill()
    this.exocom.close(done)
  })


  AfterAll(function(_scenario, done) {
    getDb( (db) => {
      db.collection('todos').drop()
      db.close(function(err, result){
        if (err) { throw new Error(err) }
        done()
      })
    })
  })

})
