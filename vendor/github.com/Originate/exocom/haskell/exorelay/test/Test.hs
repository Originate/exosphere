module Main where

import Network.Exocom
import Control.Concurrent
import Control.Concurrent.MVar
import System.Process
import System.Exit
import Data.Aeson
import Data.ByteString as B hiding (putStrLn)
import qualified Data.ByteString.Char8 as SB
import qualified Data.ByteString.Lazy as LB


echoHandler :: MVar Value -> Value -> IO ()
echoHandler = putMVar

main :: IO ()
main = do
  exo <- newExoRelay 4100 "exorelay-hs" 4001 (Just handleError)
  didRoundTrip <- roundTrip exo
  didSendReply <- testSendReply exo
  didListenReply <- testReply exo
  let results = [didRoundTrip, didSendReply, didListenReply]
  if and results then exitSuccess else exitFailure


roundTrip :: ExoRelay -> IO Bool
roundTrip exo = do
  ctrlVar <- newEmptyMVar
  registerHandler exo "hello" (echoHandler ctrlVar)
  sendMsg exo "hello" (toJSON "payload")
  res <- readMVar ctrlVar
  return $ res == (toJSON "payload") where

testSendReply :: ExoRelay -> IO Bool
testSendReply exo = do
  ctrlVar <- newEmptyMVar
  sendMsgWithReply exo "reply" (toJSON "reply Payload") (\val -> do
    if val == (toJSON "reply Payload") then putMVar ctrlVar True
    else putMVar ctrlVar False
    )
  result <- readMVar ctrlVar
  return result

testReply :: ExoRelay -> IO Bool
testReply exo = do
  ctrlVar <- newEmptyMVar
  registerHandlerWithReply exo "needReply" (\val -> do
    registerHandler exo "listenReply" (\val -> putMVar ctrlVar (val == (toJSON "payload")))
    return ("listenReply", val)
    )
  result <- readMVar ctrlVar
  return result


handleError :: String -> IO ()
handleError str = do
  Prelude.putStrLn $ "Error found: " ++ str
