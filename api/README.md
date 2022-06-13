# api

`Application Program Interfaces`. This folder contains the interface implementations for the `core` and `system`
layers to communicate with the outside world, i.e. http, rpc, binary, sockets, etc. 

The `core` & `system` layers will define interfaces about how they ought to be used 
and your api layer will build implementations for them.

Decoupling the API from the system logic allows us to change the system and domain 
without affecting the API and vice versa. 

This is also known as the `interfaces` layer according to the `clean architecture` design philosophy.

Try 

```go
TRY(adding another database implementation alongside postgres)
TRY(adding a Socket or RPC server alongside the gin http)
TRY(adding another http server alongside gin)
```
