package exorelay

import akka.actor._
import akka.testkit._
import com.softwaremill.tagging._
import ConnectionManager.{Connection, Offline, Uninitialized}

class ConnectionManagerSpec extends CommonAkkaTest("ConnectionManagerSpec") {

  val config = Config("localhost", 4100, "test-service")

  val handler = TestProbe()
  val listener = system.actorOf(MockListener.props(config, handler.ref.taggedWith[Handler]))

  "A ConnectionManager" should "Go online and respond when it receives the Connect command " in {
    listener ! ConnectCmd
    handler.expectMsg(ConnectedEvt)
  }

  it should "Go offline when it receives the Disconnect command" in {
    listener ! DisconnectCmd
    handler.expectMsg(DisconnectedEvt)
  }
}

class MockListener(val config: Config, val out: ActorRef @@ Handler)
                  (implicit system: ActorSystem)
  extends ConnectionManager {

  override def connect(): State = {
    val testSocket = TestProbe().ref
    goto(ConnectionManager.Online) using Connection(testSocket)
  }

  override def onlineState: StateFunction = {
    case Event(e, Connection(socket)) =>
      socket ! PoisonPill
      goto(Offline) using Uninitialized
  }
}

object MockListener {
  def props(config: Config, out: ActorRef @@ Handler)(implicit system: ActorSystem) =
    Props(new MockListener(config, out))
}