package exorelay

import akka.actor._
import com.softwaremill.tagging.@@

// ConnectionManager handles connecting and disconnecting an active connection (e.g. socket)
trait ConnectionManager
  extends FSM[ConnectionManager.ConnectionState, ConnectionManager.ConnectionData] {
  import ConnectionManager._

  val config: Config

  val out: ActorRef @@ Handler

  // Create a new connection
  // @return the next connection state
  def connect(): State

  // We need a way to handle commands while online
  def onlineState: StateFunction

  // Initial state
  startWith(Offline, Uninitialized)

  // When Connection is offline
  // and it receives a connect command
  // it should go online with a new connection
  when(Offline){
    case Event(ConnectCmd, _) => connect()
  }

  when(Busy){
    case Event(e, _) =>
      log.error(s"Received message $e from $sender while busy")
      stay
  }

  when(Online)(onlineState)

  onTransition{
    // Send "Connected" event
    case Offline -> Online =>
      nextStateData match {
        case Connection(socket) =>
          out ! ConnectedEvt
        case _ =>
          log.error("Missing active connection")
      }

    // Send "Disconnected" event
    case Online -> Offline =>
      out ! DisconnectedEvt
  }

  whenUnhandled{
    case Event(message: ExoMessage, Uninitialized) =>
      log.error(s"Unable to send message $message while not online")
      stay
    case Event(e, _) =>
      log.error(s"Unhandled message $e received from $sender")
      stay
  }
}


object ConnectionManager {
  // State
  sealed trait ConnectionState
  case object Offline extends ConnectionState
  case object Online extends ConnectionState
  case object Busy extends ConnectionState

  // Data
  sealed trait ConnectionData
  case object Uninitialized extends ConnectionData
  final case class Connection(socket: ActorRef) extends ConnectionData
}