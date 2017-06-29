module Network.Exocom.Internal.ExoRelay where

import Data.ByteString as B
import qualified Data.HashMap as HM
import Control.Concurrent.Chan
import Control.Concurrent.MVar
import Data.Aeson

data MessageHandler = NoReply ListenHandler | Reply ReplyHandler

type ListenHandler = Value -> IO ()

type ReplyHandler = Value -> IO (String, Value)

data ExoRelay = ExoRelay {
  port :: Int,
  serviceName :: String,
  sendChan :: Chan B.ByteString,
  errorChan :: Chan String,
  receiveHandlers :: MVar (HM.Map String MessageHandler),
  errorHandler :: Maybe (String -> IO ())
}


-- Takes in an exorelay and a ByteString which corresponds to the command type
-- When a message comes for that message type it will execute the given handler
-- in a separate thread (WARNING: PLEASE DON'T PERFORM ANY NON-THREAD-SAFE OPERATIONS IN THE HANDLER)
registerHandler :: ExoRelay -> String -> (Value -> IO ()) -> IO ()
registerHandler exo command func = do
  handlers <- takeMVar $ receiveHandlers exo
  let newHandlers = HM.insert command (NoReply func) handlers
  putMVar (receiveHandlers exo) newHandlers

registerHandlerWithReply :: ExoRelay -> String -> (Value -> IO (String, Value)) -> IO ()
registerHandlerWithReply exo command func = do
  handlers <- takeMVar $ receiveHandlers exo
  let newHandlers = HM.insert command (Reply func) handlers
  putMVar (receiveHandlers exo) newHandlers


unregisterHandler :: ExoRelay -> String -> IO ()
unregisterHandler exo command = do
  handlers <- takeMVar $ receiveHandlers exo
  let newHandlers = HM.delete command handlers
  putMVar (receiveHandlers exo) newHandlers
