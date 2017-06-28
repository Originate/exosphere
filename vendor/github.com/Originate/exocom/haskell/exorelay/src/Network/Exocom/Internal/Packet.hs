{-# LANGUAGE OverloadedStrings #-}
{-# LANGUAGE DeriveGeneric #-}
module Network.Exocom.Internal.Packet where

import Data.Aeson
import Data.Aeson.Types
import Data.Maybe
import GHC.Generics


data SendPacket = SendPacket {
  name :: String,
  sender :: Maybe String,
  msgId :: String,
  payload :: Value,
  responseTo :: Maybe String
} deriving (Generic, Show)

instance ToJSON SendPacket where
  toJSON packet = object keyListFinal where
    keyList = [
      "name" .= toJSON (name packet),
      "id" .= toJSON (msgId packet),
      "payload" .= payload packet]
    keyList1
      | isJust (responseTo packet) = "responseTo" .= (toJSON (fromJust (responseTo packet))) : keyList
      | otherwise = keyList
    keyListFinal
      | isJust (sender packet) = "sender" .= (toJSON (fromJust (sender packet))) : keyList1
      | otherwise = keyList1

instance FromJSON SendPacket where
  parseJSON (Object v) = do
    nameString <- v .: "name"
    senderString <- v .:? "sender"
    msgIdString <- v .: "id"
    payloadString <- v .: "payload" :: Parser (Value)
    responseToString <- v .:? "responseTo" :: Parser (Maybe String)
    return $ SendPacket nameString senderString msgIdString payloadString responseToString
