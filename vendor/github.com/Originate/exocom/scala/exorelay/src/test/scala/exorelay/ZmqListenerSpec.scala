package exorelay

import akka.actor.{ActorRef, ActorSystem}
import akka.pattern._
import akka.testkit._
import org.scalatest._
import zeromq.ZeroMQExtension
import ZmqListener.{Disconnect, Offline, Online}
import akka.util.Timeout
import org.zeromq.ZMQException

import scala.concurrent.duration._

class ZmqListenerSpec
  extends TestKit(ActorSystem("ZmqListenerSpec"))
  with ImplicitSender
  with FlatSpecLike
  with Matchers
  with BeforeAndAfterAll {

  type TestZmqListener = TestFSMRef[ZmqListener.State, ZmqListener.Data, ZmqListener]

  val MAX_DURATION = 10.seconds

  implicit val executor = system.dispatcher

  val zmq = ZeroMQExtension(system)

  def withTestFSM(testCode: (TestZmqListener, TestProbe) => Any) = {
    val probe = TestProbe()
    val zmqListener = TestFSMRef(new ZmqListener(out = probe.ref, zmq))
    testCode(zmqListener, probe)
  }

  override def afterAll = {
    TestKit.shutdownActorSystem(system)
  }

  "A ZmqListener" should "Manage a connection to an endpoint for a given port" in {

      val fsm = system.actorOf(ZmqListener.props(out = self, zmq))

      fsm ! ZmqListener.Connect(4100)

      expectMsgPF(max = MAX_DURATION){
        case msg @ ZmqListener.Connected(_) => ()
      }

      fsm ! ZmqListener.Disconnect
      expectMsg(ZmqListener.Disconnected)
  }

  it should "Reply with an error if the connection fails" in withTestFSM {
    (fsm, probe) =>

      fsm.setState(Offline)
      fsm.setState(ZmqListener.Failure)

      probe.expectMsgPF(max = MAX_DURATION){
        case msg @ ZmqListener.ConnectionFailed(_) => ()
      }
  }
}
