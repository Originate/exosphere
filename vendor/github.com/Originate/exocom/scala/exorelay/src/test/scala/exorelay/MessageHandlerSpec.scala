package exorelay

import akka.actor.{ActorSystem, Props}
import akka.pattern._
import akka.testkit._
import org.scalatest._
import zeromq._

class MessageHandlerSpec
  extends TestKit(ActorSystem("ConnectionManagerSpec"))
    with ImplicitSender
    with FlatSpecLike
    with Matchers
    with BeforeAndAfterAll {

  def withHandler(testCode: TestActorRef[MessageHandler] => Any) = {
    val handler = TestActorRef[MessageHandler]
    testCode(handler)
  }

  override def afterAll() = {
    TestKit.shutdownActorSystem(system)
  }

  "A MessageHandler" should "Add a new handler by sending the RegisterHandler command" in withHandler {
    handler =>

      val addCommand = MessageHandler.AddHandler("test", _ => ())

      handler ! addCommand

      expectMsg(true)
  }

  it should "Remove an existing handler by sending the RemoveHander command" in withHandler {
    handler =>

      val addCommand = MessageHandler.AddHandler("test.good", _ => ())
      handler ! addCommand

      val removeCommand = MessageHandler.RemoveHandler("test.good")
      handler ! removeCommand

      expectMsg(true) // Added handler
      expectMsg(true) // Removed handler

      val badRemoveCommand = MessageHandler.RemoveHandler("test.bad")
      handler ! badRemoveCommand

      expectMsg(false) // Failed to remove nonexisting handler
  }

  it should "Receive incoming messages from Exosphere" in withHandler {
    handler =>

      val checkGreeting: String => Unit =
        greeting => greeting shouldBe "hello"

      val addHandler = MessageHandler.AddHandler("test.greeting", checkGreeting)
      handler ! addHandler

      expectMsg(true) // Added handler

      val message = ExoMessage(Utils.uuid, "test.greeting", "hello")
      val pickled =  upickle.default.write(message)

      handler ! pickled
      expectNoMsg()
  }

  it should "Throw an exception if it received an invalid message" in withHandler {
    handler =>

      intercept[upickle.Invalid.Json]{
        handler.receive("bad-message") // throw InvalidJson exception
      }
  }
}
