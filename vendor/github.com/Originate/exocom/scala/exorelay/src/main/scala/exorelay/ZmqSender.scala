package exorelay

import akka.actor._
import akka.util.ByteString
import com.softwaremill.tagging.@@
import exorelay.ConnectionManager.{Connection, Offline, Uninitialized}
import zeromq.{SocketType, ZeroMQExtension}

/**
  * Push messages over a zmq socket connection
  * @param zmq Used to create new ZMQ sockets
  * @param config set the host and port of the connection endpoint
  * @param out Handler to receive events
  */
class ZmqSender(zmq: ZeroMQExtension, val config: Config, val out: ActorRef @@ Handler)
  extends ConnectionManager {

  private def getEndpoint(host: String, port: Int) =
    f"tcp://$host%s:$port%d"

  def zmqMessage(exoMessage: ExoMessage): zeromq.Message = {
    val pickled = upickle.default.write(exoMessage)
    zeromq.Message(ByteString(pickled))
  }

  override def connect(): State = {
    val endpoint = getEndpoint(config.exocomHost, config.exocomPort)
    val pushSocket = zmq.newSocket(SocketType.Push, zeromq.Connect(endpoint))

    goto(ConnectionManager.Online) using ConnectionManager.Connection(pushSocket)
  }

  override def onlineState: StateFunction = {
    case Event(DisconnectCmd, Connection(socket)) =>
      socket ! PoisonPill
      goto(Offline) using Uninitialized

    case Event(message: zeromq.Message, Connection(socket)) =>
      socket ! message
      out ! MessageSentEvt
      stay

    case Event(message: ExoMessage, Connection(socket)) =>
      self ! zmqMessage(message)
      stay
  }
}

object ZmqSender {
  def props(zmq: ZeroMQExtension, config: Config, out: ActorRef @@ Handler) = Props(new ZmqSender(zmq, config, out))
}