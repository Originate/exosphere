package exorelay

/**
  * Configuration options to connect to the Exosphere service
  *
  * @param exocomHost Interface used by the sender
  * @param exocomPort Port used by the sender and receiver
  * @param serviceName Name of the service used by the sender
  */
case class Config(exocomHost: String,
                  exocomPort: Int,
                  serviceName: String)
