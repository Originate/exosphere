package exorelay

import akka.actor.{ActorRef, ActorSystem}
import akka.testkit._
import org.scalatest._

object EventBusImplSpec {
  val TEST_CHANNEL = "test-channel"
}

class EventBusImplSpec
  extends TestKit(ActorSystem("EventBusImplSpec"))
  with FlatSpecLike
  with Matchers
  with BeforeAndAfterAll {
  import EventBusImplSpec._

  override def afterAll = {
    TestKit.shutdownActorSystem(system)
  }

  def withEventBus(testCode: EventBusImpl => Any) = {
    val bus = new EventBusImpl

    testCode(bus)
  }

  def withEventBusAndSubscriber(testCode: (EventBusImpl, TestProbe) => Any) = {
    val probe = TestProbe()
    val bus = new EventBusImpl

    bus.subscribe(probe.ref, TEST_CHANNEL)

    testCode(bus, probe)
  }

  "An EventBusImpl" should "be able to subscribe an actor" in withEventBus {
    eventBus =>
      val ref = TestProbe().ref

      assert(eventBus.subscribe(ref, TEST_CHANNEL))
  }

  it should "not be able to subscribe an existing subscriber" in withEventBusAndSubscriber {
    (eventBus, subscriber) =>

      eventBus.subscribe(subscriber.ref, TEST_CHANNEL) shouldBe false
  }

  it should "be able to unsubcribe an existing subscriber" in withEventBusAndSubscriber {
    (eventBus, subscriber) =>

      eventBus.unsubscribe(subscriber.ref, TEST_CHANNEL) shouldBe true
  }

  it should "not be able to unsubscribe a non-existing subscriber" in withEventBus {
    eventBus =>
      val ref = TestProbe().ref

      eventBus.unsubscribe(ref, TEST_CHANNEL) shouldBe false
  }

  it should "publish to the subscribers in a channel" in withEventBusAndSubscriber {
    (eventBus, subscriber) =>
      val payload = "hello"
      val event = EventMessage(TEST_CHANNEL, payload)

      eventBus.publish(event)

      subscriber.expectMsg(payload)
  }
}