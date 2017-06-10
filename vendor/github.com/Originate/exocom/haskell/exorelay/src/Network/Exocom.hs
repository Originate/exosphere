{-# LANGUAGE OverloadedStrings #-}
module Network.Exocom (
  ExoRelay,
  newExoRelay,
  sendMsg,
  sendMsgWithReply,
  registerHandler,
  registerHandlerWithReply,
)
where

import System.ZMQ4
import Control.Concurrent.MVar
import qualified Data.HashMap as HM
import Control.Concurrent.Chan
import Control.Concurrent
import Network.Exocom.Internal.Packet
import Network.Exocom.Internal.ExoRelay
import Network.Exocom.Internal.Sender
import Network.Exocom.Internal.Listener
import Network.Exocom.Internal.Error
import Data.Maybe
import Data.Aeson


newExoRelay :: Int -> String -> Int -> Maybe (String -> IO ()) -> IO ExoRelay
newExoRelay portNum service listenerPort errHandler = do
  resetError
  let handlerMap = HM.empty
  handlerMapLock <- newMVar handlerMap -- newMVar :: IO (MVar HashMap)
  sendchan <- newChan
  errChan <- newChan
  newContext <- context
  oSock <- socket newContext Push
  iSock <- socket newContext Pull
  let exo = ExoRelay portNum service sendchan errChan handlerMapLock errHandler
  registerHandlerWithReply exo "__status" statusHandler
  errorThread exo
  sendErr <- senderThread exo oSock
  listenErr <- listenerThread exo iSock listenerPort
  return exo

statusHandler :: Value -> IO (String, Value)
statusHandler _ = return (cmd, Data.Aeson.Null) where
  cmd = "__status-ok"
  emptyJson = Data.Aeson.Null
