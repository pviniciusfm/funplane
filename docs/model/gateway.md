# Gateway

A Fanplane **Gateway** simplifies the utilization of envoy core features such as monitoring, circuit breaking and
service to service connectivity.

Fanplane is able to inflate this template with meaningful defaults hiding the complexity from the final users.

A gateway is composed by a listener and a set of routes.
The routes varies according with the route type.
Currently, Fanplane supports the followint RouteTypes:

| Attribute     | Description                                                       |
| ------------- |:-------------------------------------------------:                |
| listener      | Define listener options see [listener](#listener)                 |
| routes        | Defines the routes exposed by the sidecar                         |
| selector      | Defines a specific match of sidecars (used in special deployments |

**Example**:
```yaml
apiVersion: fanplane.io/v1alpha1
kind: Gateway
metadata:
  name: falcon
spec:
  selector:
    version: v1

  listener:
    port: 443
    protocol: HTTPS

  routes:
  - type: CONSUL_DNS
    prefix: /iris/settings/
    service:
      name: ural2
      tag: rc0
    rewrite:
      uri: /promotion/v1/
    connectionPool:
      timeout: 5s
    retryPolicy:
      numRetries: 2
      perTryTimeout: 1s

  - type: CONSUL_DNS
    prefix: /graphql
    service:
      name: data-services
      tag: rc
    retryPolicy:
      numRetries: 2
      perTryTimeout: 1s
    tlsContext:
      sni: rc-data-services.dev.frgcloud.com

  - type: CONSUL_DNS
    prefix: /graphql
    service:
      name: data-services
      tag: rc
    retryPolicy:
      numRetries: 2
      perTryTimeout: 1s

  - type: DNS
    prefix: /flex/v1/keys
    dns:
      fqdn: api.cybersource.com
    retryPolicy:
      numRetries: 2
      perTryTimeout: 500ms

```


## Listener

**Example**:
```
listener:
    port: 443
    protocol: HTTPS
```

The listener entity define which port and protocol will be used to listen
requests.

| Value      | Description                                                                                  |
| ---------- |:----------------------------------------------------------------------------:                |
| port | Defines which port the side car will use to host the routes defined in gateway entity    |
| protocol | Defines if wich type of tcp protoco it listens (HTTP, HTTPS)                                 |


## Routes

| Value      | Description                                                                                  |
| ---------- |:----------------------------------------------------------------------------:                |
| type | Defines a [RouteType](#route-type)    |
| prefix | Defines a path utilized for route traffic                                 |
| connectionPool | Defines all configurations required to tweak how envoy manages the upstream connections                                 |
| retryPolicy | Defines a jittered backoff exponential parameters                                 |


## RouteTypes

**Example**:

| Value      | Description                                                                                  |
| ---------- |:----------------------------------------------------------------------------:                |
| CONSUL_DNS | Defines a route to consul through strict dns in envoy see [Consul DNS](#consul-dns-route)    |
| CONSUL_SDS | Defines a route to consul through eds (Not implemented yet!)                                 |
| DNS        | Defines an route through normal DNS. Most used for circuiting breaking external access.      |


## Consul DNS Route

**Example**:
```
  routes:
  - type: CONSUL_DNS
    prefix: /iris/settings/
    service:
      name: ural2
      tag: rc0
    rewrite:
      uri: /promotion/v1/
    connectionPool:
      timeout: 5s
    retryPolicy:
      numRetries: 2
      perTryTimeout: 1s
```

A consul dns route provides a simplified way to define the endpoint of the route that is behind consul service registry.

| Value      | Description                                                                                  |
| ---------- |:----------------------------------------------------------------------------:                |
| name | Routes traffic to desired service in consul    |
| tag | Defines any specified tags, if this is not used you assume that any registered service with name will be used in LB                                 |


## DNS Route

**Example**:
```
  routes:
  - type: CONSUL_DNS
    prefix: /iris/settings/
    service:
      name: ural2
      tag: rc0
    rewrite:
      uri: /promotion/v1/
    connectionPool:
      timeout: 5s
    retryPolicy:
      numRetries: 2
      perTryTimeout: 1s
```

DNS Route provides a way to route requests to external endpoints leveraging all circuiting breaking features.

## Retry Policies

| Value      | Description                                                                                  |
| ---------- |:----------------------------------------------------------------------------:                |
| numRetries | Specifies the number of attempts of receiving a valid HTTP response    |
| perTryTimeout | Defines the exponential backoff configuration                                 |

## Connection Pool

| Value      | Description                                                                                  |
| ---------- |:----------------------------------------------------------------------------:                |
| timeout | Define which timeout will be used for the connections    |
| maxConnectionPool | Defines how many persistent connections will be used for upstream connection                                 |
