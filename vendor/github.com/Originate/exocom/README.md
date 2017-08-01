<img src="documentation/logo.png" width="568" height="111" alt="logo">

[![Build Status](https://travis-ci.org/Originate/exocom.svg?branch=master)](https://travis-ci.org/Originate/exocom)

_ExoCom is a communication platform for micro-service based AI-native application ecosystems._

__AI-native:__
ExoCom provides infrastructure and conventions
for intelligent, high-level interactions between different micro-services:
- _cooperative:_ rather than telling each other what to do via remote procedure calls,
  ExoCom clients broadcast insights they gain about the domain,
  trusting that other services pay attention and do the correct things in return
- _articulate:_ messages describe domain-level events,
  in the language of the domain,
  based on a well-defined, shared, semantic
  [understanding](documentation/taxonomy.md) of what each message means
- _emergent:_ Over time, the flow of messages going over ExoCom
  forms a real-time stream of awareness
  about what is going on in the application domain.
  It contains all activities as well as the entire application state
  in an easily understandable and consumable stream format.
- _evolving:_ Tapping into this information stream,
  AI, data science, and operate components of the application
  can build up a higher-level understanding of what is going on in the application's domain.
  These insights and conclusions are shared back onto the bus,
  and can be used to adapt the application's behavior and appearance in intelligent ways.
  This seamless interaction and collaboration
  of traditional application logic and AI makes them fuse together into
  hybrid AI-native application architectures.

__optimized for heterogenous micro-service architectures__
- very low latency thanks to complete in-memory processing
- built-in [security](documentation/security.md) best practices
- based on open standards like [websockets](https://tools.ietf.org/html/rfc6455)
- [message translation](documentation/message_translation.md)
  for improved reusability of micro-services<sup>*</sup>
- _micro-service queue management:_ ExoCom does not simply shovel messages around.
  It understands that messages are sent as part of an over-arching workflow
  that is implemented via micro-services.
  After the workflow has timed out,
  ExoCom discards messages from all queues<sup>*</sup>
- [omnipresent](documentation/frontend_bridge.md): connects back-end, front-end,
  mobile, and IoT services into a real-time connectivity layer<sup>*</sup>

__ultra effective__
- client SDKs optimized for developer ergonomics and efficiency
  available for many popular languages<sup>*</sup>
- allows to develop services as lambda functions
- rapid _proto-duction_ tool:
  quickly put together a prototype,
  then efficiently and seamlessly evolve it
  into a production-grade AI-native cloud application
- real-time inspection and tracing<sup>*</sup>


## Related projects
- [Apache Kafka](https://kafka.apache.org)
- [NATS.io](http://nats.io)
- [SenecaJS](http://senecajs.org)
- [NSQ](http://nsq.io)


## Development

See our [developer guidelines](CONTRIBUTING.md)
