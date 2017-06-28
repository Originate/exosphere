module Network.Exocom.Internal.Listener where

import Network.Exocom.Internal.ExoRelay
import Network.Exocom.Internal.Packet
import Network.Exocom.Internal.Sender
import Network.Exocom.Internal.Error
import Control.Concurrent.MVar
import Control.Concurrent
import Control.Monad
import System.ZMQ4
import qualified Data.HashMap as HM
import Data.Aeson
import Data.Either
import Data.Maybe


-- Listener Functions


listenerThread :: ExoRelay -> Socket Pull -> Int -> IO ()
listenerThread exo sock listenPort = do
  let address = "tcp://*:" ++ (show (listenPort))
  bind sock address
  err <- getError
  case err of
    Nothing -> do
      forkIO $ waitAndRecv exo sock
      return ()
    Just errmsg -> emitError exo errmsg

waitAndRecv :: ExoRelay -> Socket Pull -> IO ()
waitAndRecv exo sock = do
  contents <- receive sock
  err <- getError
  case err of
    Nothing -> do
      let eitherPacket = eitherDecodeStrict contents :: Either String SendPacket
      when (isRight eitherPacket) $ do
        let packet = extract eitherPacket
        handlerMaybe <- getListenerForCommand exo packet
        when (isJust handlerMaybe) $ do
          let handler = fromJust handlerMaybe
          case handler of
            NoReply hand -> forkIO $ hand (payload packet)
            Reply hand -> forkIO $ do
              (cmd, repl) <- hand (payload packet)
              sendMsgReply exo cmd repl (msgId packet)
          return ()
    Just errmsg -> emitError exo errmsg
  waitAndRecv exo sock


-- helper function to extract right from an either
extract :: Either a b -> b
extract (Right x) = x


getListenerForCommand :: ExoRelay -> SendPacket -> IO (Maybe MessageHandler)
getListenerForCommand exo packet = do
  listeners <- readMVar (receiveHandlers exo)
  if isJust (responseTo packet) then do
    let response = fromJust $ responseTo packet
    let hand = HM.lookup response listeners
    case hand of
      Nothing -> do
        let handCmd = HM.lookup (name packet) listeners
        case handCmd of
          Nothing -> return Nothing
          Just handler -> return $ Just handler
      Just handler -> return $ Just handler
  else do
    let nameCmd = (name packet)
    case HM.lookup nameCmd listeners of
      Nothing -> return Nothing
      Just handler -> return $ Just handler
