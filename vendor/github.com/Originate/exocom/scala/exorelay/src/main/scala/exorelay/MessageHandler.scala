package exorelay

import akka.actor._
import upickle.default._

import scala.util._

/**
  * Message handling for incoming messages from Exosphere
  */
class MessageHandler
  extends Actor
    with ActorLogging
    with HandlerRegistry[String, Unit] {
  import MessageHandler._

  override def receive = {
    case AddHandler(event, callback) =>
      sender ! addHandler(event, callback)

    case RemoveHandler(event) =>
      sender ! removeHandler(event)

    case HasHandler(event) =>
      sender ! hasHandler(event)

    case message: String =>
      val exoMessage = upickle.default.read[ExoMessage](message)
      log.debug(s"Received Message $exoMessage from $sender")

      getHandler(exoMessage.messageName)
        .foreach(_.apply(exoMessage.payload))

    case e =>
      log.error(s"Unhandled message $e received from $sender in MessageHandler")
  }
}

object MessageHandler {

  case class AddHandler(event: String, callback: String => Unit)
  case class RemoveHandler(event: String)
  case class HasHandler(event: String)

  def props(): Props = Props(new MessageHandler)
}