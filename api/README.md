# api

`Application Program Interfaces`. This folder contains the interface implementations for the `core` and `system`
layers to communicate with the outside world, i.e. http, rpc, binary, sockets, etc.

The `core` & `system` layers will define interfaces about how they ought to be used and your api layer will build
implementations for them.

Decoupling the API from the system logic allows us to change the system and domain without affecting the API and vice
versa.

It is also an important component of the API layer to define a system-level interface that it will use. Why does the API
need to define the interface of the service that it will be using? Why don't we define the system-interactor's interface
with the interactor struct? Don't we already define a service-level interface? No. We define Repository-Interfaces at
the domain and system layers and at the system layer we define objects/structs that are our core implementations; these
take as an argument our Repository Interfaces.

But here is where it gets good: Our APIs do not need to implement every function our Interactor has to offer. Our http
route controllers might want to, but let's say we want to add a Kafka message queue that just cares about adding Txp to
a Skill. It [, the Kafka Interactor/Service]
shouldn't need to implement each function the System-Interactor does. It can define an interface with just the one
function it needs from the System-Interactor and let that Interactor fulfill its contract.

tldr: let each API define just the functions it needs to do its job and let the System-Interactor implicitly fulfill
that interface. Pass the System-Interactor to the Api-Interactor via the Api-Interactor's System Interface.

In our whole application we have Interactors which are implementations, structs, objects, nouns. They do the work. We
have them on every level. Then we use Interfaces to bridge and couple the Interactors together, but in a way that allows
easy decoupling.

This is also known as the `interfaces` layer according to the `clean architecture` design philosophy.

Try

```go
TRY(adding another database implementation alongside postgres)
TRY(adding a Socket or RPC server alongside the gin http)
TRY(adding another http server alongside gin)
```
