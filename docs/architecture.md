# Architecture

## Components

1. **adapter** - Module responsible for the translation of Fanplane Objects to Envoy entities
1. **Server** - Server module contains all code required to Fanplane's GRPC server and middleware
1. **Registry** - Is responsible to fetch and validate the configuration storage for Fanplane
1. **Apis** - Defines Kubernetes CRDs and Fanplane Entities

## Registry Controller

![Controller Example](https://github.com/kubernetes/sample-controller/raw/master/docs/images/client-go-controller-interaction.jpeg)

## Deployment types

1. [Envoy Deployment types](https://www.envoyproxy.io/docs/envoy/latest/intro/deployment_types/deployment_types)
1. [service to service](https://www.envoyproxy.io/docs/envoy/latest/intro/deployment_types/service_to_service#service-to-service-egress-listener)

