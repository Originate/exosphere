package exorelay

import akka.actor._
import zeromq._

import scala.concurrent.ExecutionContextExecutor

/**
  * ZmqListener connects to a ZMQ endpoint to listen for messages pushed by Exosphere
  * @param out Recipient of outgoing event messages
  * @param zmq ZmqExtension to create new socket connections
  */
class ZmqListener(out: ActorRef, zmq: ZeroMQExtension)
  extends FSM[ZmqListener.State, ZmqListener.Data]
    with ActorLogging {
  import ZmqListener._

  startWith(Offline, Uninitialized)

  // When ZmqListener is offline
  // it can be requested to connect
  // and begin listening on a given port
  when(Offline){
    case Event(Connect(port), _) =>
      val endpoint = s"tcp://*:$port"
      val zmqSocket = zmq.newSocket(SocketType.Pull, Bind(endpoint), Listener(out))

      goto(Online) using ConnectionData(zmqSocket)
  }

  // When ZmqListener is online
  // it can be requested to close
  // the current socket and disconnect
  when(Online){
    case Event(Disconnect, ConnectionData(socket)) =>
      socket ! PoisonPill

      goto(Offline) using Uninitialized
  }

  when(Failure){
    // Reply with a failure if failed to connect
    case _ =>
      out ! ConnectionFailed("Failed to connect")
      stay
  }

  onTransition{
    // Emit Connected status
    case Offline -> Online =>
      withConnection(nextStateData){ data =>
        out ! Connected(data.socket)
      }

    // Emit Disconnected status
    case Online -> Offline =>
      out ! Disconnected

    // Emit Failure status
    case Offline -> Failure =>
      out ! ConnectionFailed("Failed to connect")

    // Emit Failure status
    case Online -> Failure =>
      out ! ConnectionFailed("Connection failed unexpectedly")
  }

  whenUnhandled{
    case Event(e, _) =>
      log.error(s"Unhandled message $e received from $sender in zmqListener")
      stay
  }
}

object ZmqListener {

  val MAX_RETRIES = 3

  // Props
  def props(out: ActorRef, zmq: ZeroMQExtension): Props =
    Props(new ZmqListener(out, zmq))

  // State
  sealed trait State
  case object Offline extends State
  case object Online extends State
  case object Failure extends State

  // Data
  sealed trait Data
  case object Uninitialized extends Data
  final case class ConnectionData(socket: ActorRef) extends Data

  // Received Commands
  sealed trait Command
  case object Disconnect extends Command
  final case class Connect(port: Int) extends Command

  // Sent Events
  sealed trait Event
  case object Disconnected extends Event
  final case class Connected(socket: ActorRef) extends Event
  final case class ConnectionFailed(msg: String) extends Event

  // Helpers
  def withConnection(data: Data)(cb: ConnectionData => Unit) = data match {
    case connection @ ConnectionData(_) => cb(connection)
    case _ => println("Error: Missing ZmqListener active connection")
  }
}
