# Communication over a message bus

Why use a message bus for communication between services,
instead of relying on direct point-to-point communication between services?

- the order in which messages have been sent and received is clear and non-debatable.
  If services send each other messages, this can be a lot more unclear.
- you get real-time logging and analytics for free. With direct communication between services,
  a service has to send two messages: one to the target service, the other to the analytics service
- messages can be buffered and re-sent in case the target service is temporarily unavailable.
  This improves the reliability of the distributed application
- you have a clear record of what exactly happened when,
  in case there are reviews or audits
- network hops: yes, communication over a message bus requires two network hops instead of one.
  So does direct communication between services when going through a load balancers, though.
- central point of contention: yes, a message bus is a central point of contention or failure.
  So is a load balancer. Like load balancers, message buses can be built to scale horizontally.
