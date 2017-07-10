package exorelay

import akka.actor.ActorRef
import akka.event.{EventBus, LookupClassification}

case class EventMessage(channel: String, payload: Any)

class EventBusImpl extends EventBus with LookupClassification {
  override type Event = EventMessage
  override type Classifier = String
  override type Subscriber = ActorRef

  // Define a full order over the subscribers
  override protected def compareSubscribers(a: Subscriber, b: Subscriber): Int =
    a.compareTo(b)

  // Extract the classifier from the event
  override protected def classify(event: Event): Classifier =
    event.channel

  // Invoked for all subscribers registered the the events classifier
  override protected def publish(event: Event, subscriber: Subscriber) = {
    subscriber ! event.payload
  }

  // Expected max number of classifiers
  override protected def mapSize: Int = 10
}