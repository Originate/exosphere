# ExoCom FrontEnd Bridge

ExoCom does not only connect backend services.
It is an _omnipresent_ messaging platform
that connects _all_ parts of AI-native application ecosystem,
including services in various backends, in browsers,
and on mobile and IoT devices.

[[ image ]]

The backend portion of Exocom connects backend services.
The _public bridge_ maintains a number of websockets
to the various front-end clients (browsers, mobile devices).
It is aware about which user is logged in via which websocket connection,
and forwards messages that are supposed to reach the user on front-end
devices over the respective sockets.

