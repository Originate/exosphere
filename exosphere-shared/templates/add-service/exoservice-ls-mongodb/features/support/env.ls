process.env.NODE_ENV = 'test'
require! {
  'mongodb' : {MongoClient}
  'nitroglycerin' : N
}


db = null
get-db = (done) ->
  return done db if db
  MongoClient.connect "mongodb://localhost:27017/exosphere-_____serviceRole_____-test", N (mongo-db) ->
    db := mongo-db
    done db


module.exports = ->

  @set-default-timeout 1000


  @Before (_scenario, done) ->
    get-db (db) ->
      db.collection('_____modelName_____s')?.drop!
      done!

  @After ->
    @exocom?.close!
    @process?.close!


  @registerHandler 'AfterFeatures', (_event, done) ->
    get-db (db) ->
      db.collection('_____modelName_____s')?.drop!
      db.close!
      done!

