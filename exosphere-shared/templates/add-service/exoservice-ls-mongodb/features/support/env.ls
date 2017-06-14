process.env.NODE_ENV = 'test'
require! {
  'cucumber': {defineSupportCode}
  'mongodb' : {MongoClient}
  'nitroglycerin' : N
  './world': World
}


db = null
get-db = (done) ->
  return done db if db
  MongoClient.connect "mongodb://localhost:27017/exosphere-_____serviceRole_____-test", N (mongo-db) ->
    db := mongo-db
    done db


defineSupportCode ({After, Before, set-default-timeout, set-world-constructor, registerHandler}) ->

  set-default-timeout 1000
  set-world-constructor World


  Before (_scenario, done) ->
    get-db (db) ->
      db.collection('_____modelName_____s')?.drop!
      done!

  After ->
    @exocom?.close!
    @process?.close!


  registerHandler 'AfterFeatures', (_event, done) ->
    get-db (db) ->
      db.collection('_____modelName_____s')?.drop!
      db.close (err, result) ->
        | err => throw new Error err
        done!
