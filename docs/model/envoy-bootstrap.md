# EnvoyBootstrap

An EnvoyBootstrap endpoint provides a smooth transition way from static envoy configuration
to a Fanplane dynamic objects. Behind the scenes Fanplane will convert the
EnvoyBootstrap messages into CDS and LDS envoy messages.

## Bootstrap Header

In Fanplane the header object defines which sidecar the desired configuration will take effect.

```
apiVersion: networking.fanplane.io/v1alpha1
kind: EnvoyBootstrap
metadata:
  name: falcon //Defines the name of the sidecar
spec:
```

| Attribute     | Description                                                       |
| ------------- |:-------------------------------------------------:                |
| name          | Define which sidecar will be the target of the entities specified in the entity|
| spec          | Marks the beginning of the static envoy configuration|

Everything under spec is validated by Bootstrap proto message [see](https://www.envoyproxy.io/docs/envoy/latest/api-v2/bootstrap/bootstrap)


## Example configuration

```
apiVersion: networking.fanplane.io/v1alpha1
kind: EnvoyBootstrap
metadata:
  name: falcon
  tags:
    version: v1
spec:
  admin:
    access_log_path: /gateway/admin_access.log
    address:
      socket_address: { address: 127.0.0.1, port_value: 9001 }
  static_resources:
    listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 8800 }
      filter_chains:
      - filters:
        - name: envoy.http_connection_manager
          config:
            stat_prefix: ingress_http
            access_log:
            - name: envoy.file_access_log
              config:
                path: /gateway/envoy_access.log
                format: "[%START_TIME%] \"%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL%\" %RESPONSE_CODE% %RESPONSE_FLAGS% %BYTES_RECEIVED% %BYTES_SENT% %DURATION% %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% \"%REQ(X-FORWARDED-FOR)%\" \"%REQ(USER-AGENT)%\" \"%REQ(X-REQUEST-ID)%\" \"%REQ(:AUTHORITY)%\" \"%UPSTREAM_HOST%\" \"%REQ(X-FRG-RC)%\" \"%REQ(X-FRG-RI)%\"\n"
            route_config:
              name: local_route
              virtual_hosts:
              - name: local_service
                domains: ["*"]
                routes:
                - match: { prefix: "/iris/settings/" }
                  route:
                    host_rewrite: rc0.ural2.service.us-east-1.dynamic.prod.frgcloud.com
                    cluster: service_ural
                    timeout: 5s
                    retry_policy: {retry_on: 5xx, num_retries: 2, per_try_timeout: 1s}

                - match: { prefix: "/graphql" }
                  route:
                    host_rewrite: rc.data-services.service.us-east-1.dynamic.prod.frgcloud.com
                    cluster: service_ds
                    timeout: 5s
                    retry_policy: {retry_on: 5xx, num_retries: 2, per_try_timeout: 1s}

                - match: { prefix: "/v2/experience" }
                  route: { host_rewrite: rc.experience.service.us-east-1.dynamic.prod.frgcloud.com, cluster: service_experience }

                - match: { prefix: "/v2/cart" }
                  route: { host_rewrite: rc.order.service.us-east-1.dynamic.prod.frgcloud.com, cluster: service_order }

                - match: { prefix: "/v1/account" }
                  route: { host_rewrite: rc.user.service.us-east-1.dynamic.prod.frgcloud.com, cluster: service_user }

                - match: { prefix: "/v1/validate" }
                  route: { host_rewrite: rc.custom-prod.service.us-east-1.dynamic.prod.frgcloud.com, cluster: service_customprod }

                - match: { prefix: "/client/v1/token" }
                  route: { host_rewrite: rc.clientauth.service.us-east-1.dynamic.prod.frgcloud.com, cluster: service_clientauth }

                - match: { prefix: "/flex/v1/keys" }
                  route:
                    host_rewrite: api.cybersource.com
                    cluster: service_cybersource
                    timeout: 3s
                    retry_policy: {retry_on: 5xx, num_retries: 2, per_try_timeout: 0.5s}

            http_filters:
            - name: envoy.router
    clusters:
    - name: service_ural
      connect_timeout: 5s
      type: STRICT_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: rc0.ural2.service.us-east-1.dynamic.prod.frgcloud.com, port_value: 8443 }}]
      tls_context: { sni: rc0-ural2.prod.frgcloud.com }
      circuit_breakers:
        thresholds:
        - priority: DEFAULT
          max_connections: 1024
          max_pending_requests: 1024
          max_retries: 3

    - name: service_ds
      connect_timeout: 5s
      type: STRICT_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: rc.data-services.service.us-east-1.dynamic.prod.frgcloud.com, port_value: 8443 }}]
      tls_context: { sni: rc-data-services.prod.frgcloud.com }
      circuit_breakers:
        thresholds:
        - priority: DEFAULT
          max_connections: 1024
          max_pending_requests: 1024
          max_retries: 3

    - name: service_cybersource
      connect_timeout: 3s
      type: LOGICAL_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: api.cybersource.com, port_value: 8443 }}]
      tls_context: { sni: api.cybersource.com }

    - name: service_experience
      connect_timeout: 5s
      type: STRICT_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: rc.experience.service.us-east-1.dynamic.prod.frgcloud.com, port_value: 8443 }}]
      tls_context: { sni: rc.experience.service.us-east-1.dynamic.prod.frgcloud.com }

    - name: service_order
      connect_timeout: 5s
      type: STRICT_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: rc.order.service.us-east-1.dynamic.prod.frgcloud.com, port_value: 8443 }}]
      tls_context: { sni: rc.order.service.us-east-1.dynamic.prod.frgcloud.com }

    - name: service_user
      connect_timeout: 5s
      type: STRICT_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: rc.user.service.us-east-1.dynamic.prod.frgcloud.com, port_value: 8443 }}]
      tls_context: { sni: rc.user.service.us-east-1.dynamic.prod.frgcloud.com }

    - name: service_customprod
      connect_timeout: 5s
      type: STRICT_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: rc.custom-prod.service.us-east-1.dynamic.prod.frgcloud.com, port_value: 8443 }}]
      tls_context: { sni: rc.custom-prod.service.us-east-1.dynamic.prod.frgcloud.com }

    - name: service_clientauth
      connect_timeout: 5s
      type: STRICT_DNS
      dns_lookup_family: V4_ONLY
      lb_policy: ROUND_ROBIN
      hosts: [{ socket_address: { address: rc.clientauth.service.us-east-1.dynamic.prod.frgcloud.com, port_value: 8443 }}]
      tls_context: { sni: rc.clientauth.service.us-east-1.dynamic.prod.frgcloud.com }


```

