package exorelay

import akka.actor.{ActorRef, ActorSystem}
import akka.pattern._
import akka.testkit._
import org.scalatest._
import com.softwaremill.tagging._
import zeromq.ZeroMQExtension
import ZmqListener.{Disconnect, Offline, Online}
import akka.util.Timeout
import org.zeromq.ZMQException

class ZmqSenderSpec extends CommonAkkaTest("ZmqSenderSpec") {

  val zmq = ZeroMQExtension(system)

  val config = Config("localhost", 4100, "zmqsender-service")

  val probe = TestProbe()

  val sender = system.actorOf(ZmqSender.props(zmq, config, probe.ref.taggedWith[Handler]))

  "A ZmqSender" should "Connect using a new socket and reply when it's ready" in {
    sender ! ConnectCmd
    probe.expectMsg(ConnectedEvt)
  }

  it should "Send a new message and reply when it's sent" in {
    val message = ExoMessage("123","test","hello")

    sender ! message
    probe.expectMsg(MessageSentEvt)
  }
}
