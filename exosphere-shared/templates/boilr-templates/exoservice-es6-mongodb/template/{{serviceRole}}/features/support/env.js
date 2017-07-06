process.env.NODE_ENV = 'test'
const defineSupportCode = require('cucumber').defineSupportCode,
      {MongoClient} = require('mongodb'),
      N = require('nitroglycerin'),
      World = require('./world')


var db = null
const getDb = (done) => {
  if (db) return done(db)
  MongoClient.connect("mongodb://localhost:27017/exosphere-{{serviceRole}}-test", N( (mongoDb) => {
    db = mongoDb
    done(db)
  }))
}


defineSupportCode(function({Before, After, setDefaultTimeout, setWorldConstructor, registerHandler}) {

  setDefaultTimeout(1000)
  setWorldConstructor(World)


  Before( function(_scenario, done) {
    getDb( (db) => {
      db.collection('{{modelName}}s').drop(function(err) {
        // ignore errors here, since we are only cleaning up the test database
        // and it might not even exist
        done()
      })
    })
  })


  After(function() {
    this.exocom && this.exocom.close()
    this.process && this.process.close()
  })


  registerHandler('AfterFeatures', (_event, done) => {
    getDb( (db) => {
      db.collection('{{modelName}}s').drop()
      db.close(function(err, result){
        if (err) { throw new Error(err) }
        done()
      })
    })
  })

});
