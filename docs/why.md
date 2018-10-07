## Why Envoy

[Compare Envoy](https://www.envoyproxy.io/docs/envoy/latest/intro/comparison)

## Why use Fanplane?

While solid open source control planes like Istio provides a swiss-knife for developers they don't provide
a transition from non service-mesh to a service mesh architecture.

In fanatics we needed a way to better control our envoy sidecars and also
enforce best practices but we couldn't afford to migrate all our services to containers and kubernetes all it once.

Fanplane is a way to provide a smooth transition to service mesh based architecture as it offers an hybrid model where
you can optionally have fine control about all features of envoy sidecars

Fanplane provides an quick to way to adopt service-to-service egress listener envoy deployment model
[see for reference](https://www.envoyproxy.io/docs/envoy/latest/intro/deployment_types/service_to_service#service-to-service-egress-listener)