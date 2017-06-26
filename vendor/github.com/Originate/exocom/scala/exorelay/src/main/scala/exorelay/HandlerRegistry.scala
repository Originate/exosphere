package exorelay

/**
  * Generic registry to handle events
  */
trait HandlerRegistry[In, Out] {
  private val handlers = scala.collection.mutable.AnyRefMap.empty[String, In => Out]

  def getHandler(event: String) = handlers.get(event)

  def addHandler(event: String, callback: In => Out): Boolean = {
    if(handlers.contains(event)) false
      else {
      handlers += (event -> callback)
      true
    }
  }

  def removeHandler(event: String): Boolean = {
    if(handlers.contains(event)) {
      handlers -= event
      true
    }
    else false
  }

  def hasHandler(event: String): Boolean = handlers.contains(event)
}
