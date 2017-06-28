package exorelay

import akka.actor.ActorRef

/**
  * Exocom Message sent and received
  * @param id Generated identifier for the message
  * @param messageName nNme of the message sent (e.g. 'user.created')
  * @param payload Optional additional data to be sent with the message
  * @param responseTo Optional response to some message
  */
case class ExoMessage(id: String, messageName: String, payload: String, responseTo: Option[String] = None)

trait Listener
trait Sender
trait Handler

// Received commands
sealed trait ConnectionCmd
case object DisconnectCmd extends ConnectionCmd
case object ConnectCmd extends ConnectionCmd
final case class SendCmd(name: String, payload: Any) extends ConnectionCmd

// Sent events
sealed trait ConnectionEvt
case object ConnectedEvt extends ConnectionEvt
case object DisconnectedEvt extends ConnectionEvt
case object MessageSentEvt extends ConnectionEvt