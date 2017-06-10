module Network.Exocom.Internal.Sender where

import System.ZMQ4
import Network.Exocom.Internal.ExoRelay
import Network.Exocom.Internal.Error
import qualified Data.ByteString.Lazy as LB
import Control.Concurrent.Chan
import Control.Concurrent
import Data.UUID
import Data.UUID.V4
import Data.Aeson
import Network.Exocom.Internal.Packet



-- Sender Functions


senderThread :: ExoRelay -> Socket Push -> IO ()
senderThread exo sock = do
  let address = "tcp://localhost:" ++ (show (port exo))
  resetError
  connect sock address
  err <- getError
  case err of
    Nothing -> forkIO (waitAndSend exo sock) >> return ()
    Just errmsg -> emitError exo errmsg

waitAndSend :: ExoRelay -> Socket Push -> IO ()
waitAndSend exo sock = do
  toSend <- readChan $ sendChan exo
  send sock [] toSend
  err <- getError
  case err of
    Nothing -> return ()
    Just errmsg -> emitError exo errmsg
  waitAndSend exo sock


-- internal sending of msg
sendMsgGeneral :: ExoRelay -> String -> Value -> Maybe String -> IO ()
sendMsgGeneral exo command toSend respond = do
  identUUID <- nextRandom
  let ident = toString identUUID
  let packet = SendPacket command (Just (serviceName exo)) ident toSend respond
  let jsonByteString = encode packet
  writeChan (sendChan exo) (LB.toStrict jsonByteString)


-- sendMsg takes in the exorelay object, a command type and a payload and sends it
sendMsg :: ExoRelay -> String -> Value -> IO ()
sendMsg exo command toSend = sendMsgGeneral exo command toSend Nothing


-- sendMsgReply acts like sendMsg but has a last argument which is a UUID to which the message is replying to
sendMsgReply :: ExoRelay -> String -> Value -> String -> IO ()
sendMsgReply exo cmd toSend replUUID = sendMsgGeneral exo cmd toSend (Just replUUID)

sendMsgWithReply :: ExoRelay -> String -> Value -> (Value -> IO ()) -> IO ()
sendMsgWithReply exo cmd payload hand = do
  identUUID <- nextRandom
  let ident = toString identUUID
  let packet = SendPacket cmd (Just (serviceName exo)) ident payload Nothing
  let jsonByteString = encode packet
  registerHandler exo ident (\response -> unregisterHandler exo ident >> hand response)
  writeChan (sendChan exo) (LB.toStrict jsonByteString)
