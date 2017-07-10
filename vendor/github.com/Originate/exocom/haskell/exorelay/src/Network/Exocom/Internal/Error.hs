module Network.Exocom.Internal.Error (
  getError,
  resetError,
  emitError,
  errorThread
)
where

import Network.Exocom.Internal.ExoRelay
import Control.Concurrent.Chan
import Control.Concurrent
import Foreign.C.Error
import Foreign.C.Types
import System.ZMQ4.Internal.Error

-- getError retrieves a ZMQ error if any and resets errno
getError :: IO (Maybe String)
getError = do
  err <- getErrno
  if err == eOK || err == eAGAIN then return Nothing
  else do
    msg <- zmqErrnoMessage (extract err)
    resetErrno
    return $ Just msg

resetError :: IO ()
resetError = resetErrno


extract :: Errno -> CInt
extract (Errno val) = val

emitError :: ExoRelay -> String -> IO ()
emitError exo errmsg = do
  let errHandler = errorHandler exo
  case errHandler of
    Nothing -> return ()
    Just hand -> do
      let errChan = errorChan exo
      writeChan errChan errmsg

errorThread :: ExoRelay -> IO ()
errorThread exo = forkIO ( waitAndCatch exo) >> return ()

waitAndCatch :: ExoRelay -> IO ()
waitAndCatch exo = do
  errmsg <- readChan (errorChan exo)
  let msgHandler = (errorHandler exo)
  case msgHandler of
    Nothing -> return ()
    Just hand -> do
      forkIO $ hand errmsg
      return ()
  waitAndCatch exo
