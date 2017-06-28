package exorelay

import akka.actor._
import akka.pattern._
import akka.util.Timeout
import zeromq.ZeroMQExtension

import scala.concurrent.{ExecutionContextExecutor, Future}
import scala.concurrent.duration._

class ExoRelay(override val config: Config)(implicit system: ActorSystem)
  extends ExoRelayLike {

  implicit val t = Timeout(5.seconds)

  val zmq = new ZeroMQExtension(system)

  val messageHandler: ActorRef = system.actorOf(MessageHandler.props(), name = "message-handler")
  val zmqListener = system.actorOf(ZmqListener.props(messageHandler, zmq), name = "zmq-listener")

  override def listen() =
    zmqListener ! ZmqListener.Connect(config.exocomPort)

  override def close() =
    zmqListener ! ZmqListener.Disconnect

  override def addHander(eventName: String, callback: String => Unit): Future[Boolean] =
    (messageHandler ? MessageHandler.AddHandler(eventName, callback)).mapTo[Boolean]

  override def removeHandler(eventName: String): Future[Boolean] =
    (messageHandler ? MessageHandler.RemoveHandler(eventName)).mapTo[Boolean]

  override def hasHandler(eventName: String): Future[Boolean] =
    (messageHandler ? MessageHandler.HasHandler(eventName)).mapTo[Boolean]

  override def send(message: ExoMessage): Unit = {
    // Send
  }
}
