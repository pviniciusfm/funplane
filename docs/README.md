# Funplane Envoy Control Plane

Funplane provides a simple way to control envoy's sidecar configuration. It makes possible to tweak low level envoy
configuration and use xDS interface.


## Funplane Flow

The idea behind Funplane is that all configurations of the sidecars are managed by a centralized registry in VCS using GitOps.
With the model in place in the registry we could use then kubectl apply or terraform apply to
ship the desired configuration inside each envoy sidecar

![Funplane Flow](img/Funplane.png)

## Reference
Funplane uses kubernetes CustomResourceDefinitions to ship
envoy configuration into sidecars.

## Quick start
```
$ funplane -c mesh.yaml -p 18000 -l info
```

## Development
Funplane uses kubernetes crds toolset to generate clients and crd definition
For update the model use:

```
$ ./hack/update_codegen.sh
```
To generate go-client and deepcopy
