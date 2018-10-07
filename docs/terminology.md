# Fanplane Terminology

**Downstream**: A downstream host connects to Envoy, sends requests, and receives responses.

**Upstream**: An upstream host receives connections and requests from Envoy and returns responses.

**Listener**: A listener is a named network location (e.g., port, unix domain socket, etc.) that can be connected to by downstream clients. Envoy exposes one or more listeners that downstream hosts connect to.

**Cluster**: A cluster is a group of logically similar upstream hosts that Envoy connects to. Envoy discovers the members of a cluster via service discovery. It optionally determines the health of cluster members via active health checking. The cluster member that Envoy routes a request to is determined by the load balancing policy.

**Mesh**: A group of hosts that coordinate to provide a consistent network topology. In this documentation, an “Envoy mesh” is a group of Envoy proxies that form a message passing substrate for a distributed system comprised of many different services and application platforms.


## Service Mesh

A service mesh is a configurable infrastructure layer for a microservices application. It makes communication between service instances flexible, reliable, and fast. The mesh provides service discovery, load balancing, encryption, authentication and authorization, support for the circuit breaker pattern, and other capabilities.

![Fanplane Flow](https://www.datocms-assets.com/2885/1513650585-istio-pilot-consul-integration.png)
