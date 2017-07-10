package exorelay

import scala.concurrent.Future

trait ExoRelayLike {
  val config: Config

  // Open a new connection on the specified port and begin receiving messages.
  def listen(): Unit

  // Close the connection and stop listening to incoming messages.
  def close(): Unit

  // Listen to internal events: "online", "offline",
  //def on(event: String, handler: => Unit): Unit

  // Send message to Exosphere
  def send(message: ExoMessage): Unit

  // Add a new handler for incoming messages
  def addHander(eventName: String, callback: String => Unit): Future[Boolean]

  // Remove an existing handler for incoming messages
  def removeHandler(eventName: String): Future[Boolean]

  // Check if handler exists
  def hasHandler(eventName: String): Future[Boolean]
}
