# "Network"
So far this is a network framing tool for tcp sockets and event handler

Im just building this for fun

## Ideas
Using protobuf to wire the message, and using a int64 to define what type of event got sent over the wire.

* protobuf
* Event Handler based on PacketType


## Todo

* correct import of connection/packet/any.proto
* connection/handler.go set correct timeout handler
* better error handling in connection/handler.go
* better error handling in server/server.go